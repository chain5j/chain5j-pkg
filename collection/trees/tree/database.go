// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tree

import (
	"errors"
	"fmt"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/database/kvstore"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"io"
	"reflect"
	"sync"
	"time"
)

// secureKeyPrefix is the database key prefix used to store trie node preimages.
var secureKeyPrefix = []byte("secure-key-")

// secureKeyLength is the length of the above prefix + 32byte hash.
const secureKeyLength = 11 + 32

// Database is an intermediate write layer between the trie data structures and
// the disk database. The aim is to accumulate trie writes in-memory and only
// periodically flush a couple tries to disk, garbage collecting the remainder.
//
// Note, the trie Database is **not** thread safe in its mutations, but it **is**
// thread safe in providing individual, independent node access. The rationale
// behind this split design is to provide read access to RPC handlers and sync
// servers even while the trie is executing expensive garbage collection.
type Database struct {
	diskdb kvstore.Database // Persistent storage for matured trie nodes

	cleans  *fastcache.Cache           // GC friendly memory cache of clean node RLPs
	dirties map[types.Hash]*cachedNode // Data and references relationships of dirty nodes
	oldest  types.Hash                 // Oldest tracked node, flush-list head
	newest  types.Hash                 // Newest tracked node, flush-list tail

	preimages map[types.Hash][]byte // Preimages of nodes from the secure trie
	seckeybuf [secureKeyLength]byte // Ephemeral buffer for calculating preimage keys

	gctime  time.Duration     // Time spent on garbage collection since last commit
	gcnodes uint64            // Nodes garbage collected since last commit
	gcsize  types.StorageSize // Data storage garbage collected since last commit

	flushtime  time.Duration     // Time spent on data flushing since last commit
	flushnodes uint64            // Nodes flushed since last commit
	flushsize  types.StorageSize // Data storage flushed since last commit

	dirtiesSize   types.StorageSize // Storage size of the dirty node cache (exc. metadata)
	childrenSize  types.StorageSize // Storage size of the external children tracking
	preimagesSize types.StorageSize // Storage size of the preimages cache

	lock sync.RWMutex
}

// rawNode is a simple binary blob used to differentiate between collapsed trie
// nodes and already encoded RLP binary blobs (while at the same time store them
// in the same cache fields).
type rawNode []byte

func (n rawNode) cache() (hashNode, bool)   { panic("this should never end up in a live trie") }
func (n rawNode) fstring(ind string) string { panic("this should never end up in a live trie") }

// rawFullNode represents only the useful data content of a full node, with the
// caches and flags stripped out to minimize its data storage. This type honors
// the same RLP encoding as the original parent.
type rawFullNode [17]node

func (n rawFullNode) cache() (hashNode, bool)   { panic("this should never end up in a live trie") }
func (n rawFullNode) fstring(ind string) string { panic("this should never end up in a live trie") }

func (n rawFullNode) EncodeRLP(w io.Writer) error {
	var nodes [17]node

	for i, child := range n {
		if child != nil {
			nodes[i] = child
		} else {
			nodes[i] = nilValueNode
		}
	}
	return rlp.Encode(w, nodes)
}

// rawShortNode represents only the useful data content of a short node, with the
// caches and flags stripped out to minimize its data storage. This type honors
// the same RLP encoding as the original parent.
type rawShortNode struct {
	Key []byte
	Val node
}

func (n rawShortNode) cache() (hashNode, bool)   { panic("this should never end up in a live trie") }
func (n rawShortNode) fstring(ind string) string { panic("this should never end up in a live trie") }

// cachedNode is all the information we know about a single cached node in the
// memory database write layer.
type cachedNode struct {
	node node   // Cached collapsed trie node, or raw rlp data
	size uint16 // Byte size of the useful cached data

	parents  uint32                // Number of live nodes referencing this one
	children map[types.Hash]uint16 // External children referenced by this node

	flushPrev types.Hash // Previous node in the flush-list
	flushNext types.Hash // Next node in the flush-list
}

// cachedNodeSize is the raw size of a cachedNode data structure without any
// node data included. It's an approximate size, but should be a lot better
// than not counting them.
var cachedNodeSize = int(reflect.TypeOf(cachedNode{}).Size())

// cachedNodeChildrenSize is the raw size of an initialized but empty external
// reference map.
const cachedNodeChildrenSize = 48

// rlp returns the raw rlp encoded blob of the cached node, either directly from
// the cache, or by regenerating it from the collapsed node.
func (n *cachedNode) rlp() []byte {
	if node, ok := n.node.(rawNode); ok {
		return node
	}
	blob, err := rlp.EncodeToBytes(n.node)
	if err != nil {
		panic(err)
	}
	return blob
}

// obj returns the decoded and expanded trie node, either directly from the cache,
// or by regenerating it from the rlp encoded blob.
func (n *cachedNode) obj(hash types.Hash) node {
	if node, ok := n.node.(rawNode); ok {
		return mustDecodeNode(hash[:], node)
	}
	return expandNode(hash[:], n.node)
}

// childs returns all the tracked children of this node, both the implicit ones
// from inside the node as well as the explicit ones from outside the node.
func (n *cachedNode) childs() []types.Hash {
	children := make([]types.Hash, 0, 16)
	for child := range n.children {
		children = append(children, child)
	}
	if _, ok := n.node.(rawNode); !ok {
		gatherChildren(n.node, &children)
	}
	return children
}

// gatherChildren traverses the node hierarchy of a collapsed storage node and
// retrieves all the hashnode children.
func gatherChildren(n node, children *[]types.Hash) {
	switch n := n.(type) {
	case *rawShortNode:
		gatherChildren(n.Val, children)

	case rawFullNode:
		for i := 0; i < 16; i++ {
			gatherChildren(n[i], children)
		}
	case hashNode:
		*children = append(*children, types.BytesToHash(n))

	case valueNode, nil:

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

// simplifyNode traverses the hierarchy of an expanded memory node and discards
// all the internal caches, returning a node that only contains the raw data.
func simplifyNode(n node) node {
	switch n := n.(type) {
	case *shortNode:
		// Short nodes discard the flags and cascade
		return &rawShortNode{Key: n.Key, Val: simplifyNode(n.Val)}

	case *fullNode:
		// Full nodes discard the flags and cascade
		node := rawFullNode(n.Children)
		for i := 0; i < len(node); i++ {
			if node[i] != nil {
				node[i] = simplifyNode(node[i])
			}
		}
		return node

	case valueNode, hashNode, rawNode:
		return n

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

// expandNode traverses the node hierarchy of a collapsed storage node and converts
// all fields and keys into expanded memory form.
func expandNode(hash hashNode, n node) node {
	switch n := n.(type) {
	case *rawShortNode:
		// Short nodes need key and child expansion
		return &shortNode{
			Key: compactToHex(n.Key),
			Val: expandNode(nil, n.Val),
			flags: nodeFlag{
				hash: hash,
			},
		}

	case rawFullNode:
		// Full nodes need child expansion
		node := &fullNode{
			flags: nodeFlag{
				hash: hash,
			},
		}
		for i := 0; i < len(node.Children); i++ {
			if n[i] != nil {
				node.Children[i] = expandNode(nil, n[i])
			}
		}
		return node

	case valueNode, hashNode:
		return n

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

// NewDatabase creates a new trie database to store ephemeral trie content before
// its written out to disk or garbage collected. No read cache is created, so all
// data retrievals will hit the underlying disk database.
func NewDatabase(diskdb kvstore.Database) *Database {
	return NewDatabaseWithCache(diskdb, 0)
}

// NewDatabaseWithCache creates a new trie database to store ephemeral trie content
// before its written out to disk or garbage collected. It also acts as a read cache
// for nodes loaded from disk.
func NewDatabaseWithCache(diskdb kvstore.Database, cache int) *Database {
	var cleans *fastcache.Cache
	if cache > 0 {
		cleans = fastcache.New(cache * 1024 * 1024)
	}

	return &Database{
		diskdb: diskdb,
		cleans: cleans,
		dirties: map[types.Hash]*cachedNode{{}: {
			children: make(map[types.Hash]uint16),
		}},
		preimages: make(map[types.Hash][]byte),
	}
}

// DiskDB retrieves the persistent storage backing the trie database.
func (db *Database) DiskDB() kvstore.KeyValueReader {
	return db.diskdb
}

// InsertBlob writes a new reference tracked blob to the memory database if it's
// yet unknown. This method should only be used for non-trie nodes that require
// reference counting, since trie nodes are garbage collected directly through
// their embedded children.
func (db *Database) InsertBlob(hash types.Hash, blob []byte) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.insert(hash, blob, rawNode(blob))
}

// insert inserts a collapsed trie node into the memory database. This method is
// a more generic version of InsertBlob, supporting both raw blob insertions as
// well ex trie node insertions. The blob must always be specified to allow proper
// size tracking.
func (db *Database) insert(hash types.Hash, blob []byte, node node) {
	// If the node's already cached, skip
	if _, ok := db.dirties[hash]; ok {
		return
	}

	// Create the cached entry for this node
	entry := &cachedNode{
		node:      simplifyNode(node),
		size:      uint16(len(blob)),
		flushPrev: db.newest,
	}
	for _, child := range entry.childs() {
		if c := db.dirties[child]; c != nil {
			c.parents++
		}
	}
	db.dirties[hash] = entry

	// Update the flush-list endpoints
	if db.oldest == (types.Hash{}) {
		db.oldest, db.newest = hash, hash
	} else {
		db.dirties[db.newest].flushNext, db.newest = hash, hash
	}
	db.dirtiesSize += types.StorageSize(types.HashLength + entry.size)
}

// insertPreimage writes a new trie node pre-image to the memory database if it's
// yet unknown. The method will make a copy of the slice.
//
// Note, this method assumes that the database's lock is held!
func (db *Database) insertPreimage(hash types.Hash, preimage []byte) {
	if _, ok := db.preimages[hash]; ok {
		return
	}
	db.preimages[hash] = hexutil.CopyBytes(preimage)
	db.preimagesSize += types.StorageSize(types.HashLength + len(preimage))
}

// node retrieves a cached trie node from memory, or returns nil if none can be
// found in the memory cache.
func (db *Database) node(hash types.Hash) node {
	// Retrieve the node from the clean cache if available
	if db.cleans != nil {
		if enc := db.cleans.Get(nil, hash[:]); enc != nil {
			return mustDecodeNode(hash[:], enc)
		}
	}
	// Retrieve the node from the dirty cache if available
	db.lock.RLock()
	dirty := db.dirties[hash]
	db.lock.RUnlock()

	if dirty != nil {
		return dirty.obj(hash)
	}

	// Content unavailable in memory, attempt to retrieve from disk
	enc, err := db.diskdb.Get(hash[:])
	if err != nil || enc == nil {
		return nil
	}
	if db.cleans != nil {
		db.cleans.Set(hash[:], enc)
	}
	return mustDecodeNode(hash[:], enc)
}

// Node retrieves an encoded cached trie node from memory. If it cannot be found
// cached, the method queries the persistent database for the content.
func (db *Database) Node(hash types.Hash) ([]byte, error) {
	// It doens't make sense to retrieve the metaroot
	if hash == (types.Hash{}) {
		return nil, errors.New("not found")
	}
	// Retrieve the node from the clean cache if available
	if db.cleans != nil {
		if enc := db.cleans.Get(nil, hash[:]); enc != nil {
			return enc, nil
		}
	}
	// Retrieve the node from the dirty cache if available
	db.lock.RLock()
	dirty := db.dirties[hash]
	db.lock.RUnlock()

	if dirty != nil {
		return dirty.rlp(), nil
	}

	// Content unavailable in memory, attempt to retrieve from disk
	enc, err := db.diskdb.Get(hash[:])
	if err == nil && enc != nil {
		if db.cleans != nil {
			db.cleans.Set(hash[:], enc)
		}
	}
	return enc, err
}

// preimage retrieves a cached trie node pre-image from memory. If it cannot be
// found cached, the method queries the persistent database for the content.
func (db *Database) preimage(hash types.Hash) ([]byte, error) {
	// Retrieve the node from cache if available
	db.lock.RLock()
	preimage := db.preimages[hash]
	db.lock.RUnlock()

	if preimage != nil {
		return preimage, nil
	}
	// Content unavailable in memory, attempt to retrieve from disk
	return db.diskdb.Get(db.secureKey(hash[:]))
}

// secureKey returns the database key for the preimage of key, as an ephemeral
// buffer. The caller must not hold onto the return value because it will become
// invalid on the next call.
func (db *Database) secureKey(key []byte) []byte {
	buf := append(db.seckeybuf[:0], secureKeyPrefix...)
	buf = append(buf, key...)
	return buf
}

// Nodes retrieves the hashes of all the nodes cached within the memory database.
// This method is extremely expensive and should only be used to validate internal
// states in test code.
func (db *Database) Nodes() []types.Hash {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var hashes = make([]types.Hash, 0, len(db.dirties))
	for hash := range db.dirties {
		if hash != (types.Hash{}) { // Special case for "root" references/nodes
			hashes = append(hashes, hash)
		}
	}
	return hashes
}

// Reference adds a new reference from a parent node to a child node.
func (db *Database) Reference(child types.Hash, parent types.Hash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.reference(child, parent)
}

// reference is the private locked version of Reference.
func (db *Database) reference(child types.Hash, parent types.Hash) {
	// If the node does not exist, it's a node pulled from disk, skip
	node, ok := db.dirties[child]
	if !ok {
		return
	}
	// If the reference already exists, only duplicate for roots
	if db.dirties[parent].children == nil {
		db.dirties[parent].children = make(map[types.Hash]uint16)
		db.childrenSize += cachedNodeChildrenSize
	} else if _, ok = db.dirties[parent].children[child]; ok && parent != (types.Hash{}) {
		return
	}
	node.parents++
	db.dirties[parent].children[child]++
	if db.dirties[parent].children[child] == 1 {
		db.childrenSize += types.HashLength + 2 // uint16 counter
	}
}

// Dereference removes an existing reference from a root node.
func (db *Database) Dereference(root types.Hash) {
	// Sanity check to ensure that the meta-root is not removed
	if root == (types.Hash{}) {
		log().Error("Attempted to dereference the trie cache meta root")
		return
	}
	db.lock.Lock()
	defer db.lock.Unlock()

	nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	db.dereference(root, types.Hash{})

	db.gcnodes += uint64(nodes - len(db.dirties))
	db.gcsize += storage - db.dirtiesSize
	db.gctime += time.Since(start)

	log().Debug("Dereferenced trie from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
		"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)
}

// dereference is the private locked version of Dereference.
func (db *Database) dereference(child types.Hash, parent types.Hash) {
	// Dereference the parent-child
	node := db.dirties[parent]

	if node.children != nil && node.children[child] > 0 {
		node.children[child]--
		if node.children[child] == 0 {
			delete(node.children, child)
			db.childrenSize -= (types.HashLength + 2) // uint16 counter
		}
	}
	// If the child does not exist, it's a previously committed node.
	node, ok := db.dirties[child]
	if !ok {
		return
	}
	// If there are no more references to the child, delete it and cascade
	if node.parents > 0 {
		// This is a special cornercase where a node loaded from disk (i.e. not in the
		// memcache any more) gets reinjected as a new node (short node split into full,
		// then reverted into short), causing a cached node to have no parents. That is
		// no problem in itself, but don't make maxint parents out of it.
		node.parents--
	}
	if node.parents == 0 {
		// Remove the node from the flush-list
		switch child {
		case db.oldest:
			db.oldest = node.flushNext
			db.dirties[node.flushNext].flushPrev = types.Hash{}
		case db.newest:
			db.newest = node.flushPrev
			db.dirties[node.flushPrev].flushNext = types.Hash{}
		default:
			db.dirties[node.flushPrev].flushNext = node.flushNext
			db.dirties[node.flushNext].flushPrev = node.flushPrev
		}
		// Dereference all children and delete the node
		for _, hash := range node.childs() {
			db.dereference(hash, child)
		}
		delete(db.dirties, child)
		db.dirtiesSize -= types.StorageSize(types.HashLength + int(node.size))
		if node.children != nil {
			db.childrenSize -= cachedNodeChildrenSize
		}
	}
}

// Cap iteratively flushes old but still referenced trie nodes until the total
// memory usage goes below the given threshold.
//
// Note, this method is a non-synchronized mutator. It is unsafe to call this
// concurrently with other mutators.
func (db *Database) Cap(limit types.StorageSize) error {
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent storage). This is ensured
	// by only uncaching existing data when the database write finalizes.
	nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	batch := db.diskdb.NewBatch()

	// db.dirtiesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted.
	size := db.dirtiesSize + types.StorageSize((len(db.dirties)-1)*cachedNodeSize)
	size += db.childrenSize - types.StorageSize(len(db.dirties[types.Hash{}].children)*(types.HashLength+2))

	// If the preimage cache got large enough, push to disk. If it's still small
	// leave for later to deduplicate writes.
	flushPreimages := db.preimagesSize > 4*1024*1024
	if flushPreimages {
		for hash, preimage := range db.preimages {
			if err := batch.Put(db.secureKey(hash[:]), preimage); err != nil {
				log().Error("Failed to commit preimage from trie database", "err", err)
				return err
			}
			if batch.ValueSize() > kvstore.IdealBatchSize {
				if err := batch.Write(); err != nil {
					return err
				}
				batch.Reset()
			}
		}
	}
	// Keep committing nodes from the flush-list until we're below allowance
	oldest := db.oldest
	for size > limit && oldest != (types.Hash{}) {
		// Fetch the oldest referenced node and push into the batch
		node := db.dirties[oldest]
		if err := batch.Put(oldest[:], node.rlp()); err != nil {
			return err
		}
		// If we exceeded the ideal batch size, commit and reset
		if batch.ValueSize() >= kvstore.IdealBatchSize {
			if err := batch.Write(); err != nil {
				log().Error("Failed to write flush list to disk", "err", err)
				return err
			}
			batch.Reset()
		}
		// Iterate to the next flush item, or abort if the size cap was achieved. Size
		// is the total size, including the useful cached data (hash -> blob), the
		// cache item metadata, as well as external children mappings.
		size -= types.StorageSize(types.HashLength + int(node.size) + cachedNodeSize)
		if node.children != nil {
			size -= types.StorageSize(cachedNodeChildrenSize + len(node.children)*(types.HashLength+2))
		}
		oldest = node.flushNext
	}
	// Flush out any remainder data from the last batch
	if err := batch.Write(); err != nil {
		log().Error("Failed to write flush list to disk", "err", err)
		return err
	}
	// Write successful, clear out the flushed data
	db.lock.Lock()
	defer db.lock.Unlock()

	if flushPreimages {
		db.preimages = make(map[types.Hash][]byte)
		db.preimagesSize = 0
	}
	for db.oldest != oldest {
		node := db.dirties[db.oldest]
		delete(db.dirties, db.oldest)
		db.oldest = node.flushNext

		db.dirtiesSize -= types.StorageSize(types.HashLength + int(node.size))
		if node.children != nil {
			db.childrenSize -= types.StorageSize(cachedNodeChildrenSize + len(node.children)*(types.HashLength+2))
		}
	}
	if db.oldest != (types.Hash{}) {
		db.dirties[db.oldest].flushPrev = types.Hash{}
	}
	db.flushnodes += uint64(nodes - len(db.dirties))
	db.flushsize += storage - db.dirtiesSize
	db.flushtime += time.Since(start)

	log().Debug("Persisted nodes from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
		"flushnodes", db.flushnodes, "flushsize", db.flushsize, "flushtime", db.flushtime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)

	return nil
}

// Commit iterates over all the children of a particular node, writes them out
// to disk, forcefully tearing down all references in both directions. As a side
// effect, all pre-images accumulated up to this point are also written.
//
// Note, this method is a non-synchronized mutator. It is unsafe to call this
// concurrently with other mutators.
func (db *Database) Commit(node types.Hash, report bool) error {
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent storage). This is ensured
	// by only uncaching existing data when the database write finalizes.
	start := time.Now()
	batch := db.diskdb.NewBatch()

	// Move all of the accumulated preimages into a write batch
	for hash, preimage := range db.preimages {
		if err := batch.Put(db.secureKey(hash[:]), preimage); err != nil {
			log().Error("Failed to commit preimage from trie database", "err", err)
			return err
		}
		// If the batch is too large, flush to disk
		if batch.ValueSize() > kvstore.IdealBatchSize {
			if err := batch.Write(); err != nil {
				return err
			}
			batch.Reset()
		}
	}
	// Since we're going to replay trie node writes into the clean cache, flush out
	// any batched pre-images before continuing.
	if err := batch.Write(); err != nil {
		return err
	}
	batch.Reset()

	// Move the trie itself into the batch, flushing if enough data is accumulated
	nodes, storage := len(db.dirties), db.dirtiesSize

	uncacher := &cleaner{db}
	if err := db.commit(node, batch, uncacher); err != nil {
		log().Error("Failed to commit trie from trie database", "err", err)
		return err
	}
	// Trie mostly committed to disk, flush any batch leftovers
	if err := batch.Write(); err != nil {
		log().Error("Failed to write trie to disk", "err", err)
		return err
	}
	// Uncache any leftovers in the last batch
	db.lock.Lock()
	defer db.lock.Unlock()

	batch.Replay(uncacher)
	batch.Reset()

	// Reset the storage counters and bumpd metrics
	db.preimages = make(map[types.Hash][]byte)
	db.preimagesSize = 0

	logger1 := log().Info
	if !report {
		logger1 = log().Debug
	}
	logger1("Persisted trie from memory database", "nodes", nodes-len(db.dirties)+int(db.flushnodes), "size", storage-db.dirtiesSize+db.flushsize, "time", time.Since(start)+db.flushtime,
		"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)

	// Reset the garbage collection statistics
	db.gcnodes, db.gcsize, db.gctime = 0, 0, 0
	db.flushnodes, db.flushsize, db.flushtime = 0, 0, 0

	return nil
}

// commit is the private locked version of Commit.
func (db *Database) commit(hash types.Hash, batch kvstore.Batch, uncacher *cleaner) error {
	// If the node does not exist, it's a previously committed node
	node, ok := db.dirties[hash]
	if !ok {
		return nil
	}
	for _, child := range node.childs() {
		if err := db.commit(child, batch, uncacher); err != nil {
			return err
		}
	}
	if err := batch.Put(hash[:], node.rlp()); err != nil {
		return err
	}
	// If we've reached an optimal batch size, commit and start over
	if batch.ValueSize() >= kvstore.IdealBatchSize {
		if err := batch.Write(); err != nil {
			return err
		}
		db.lock.Lock()
		batch.Replay(uncacher)
		batch.Reset()
		db.lock.Unlock()
	}
	return nil
}

// cleaner is a database batch replayer that takes a batch of write operations
// and cleans up the trie database from anything written to disk.
type cleaner struct {
	db *Database
}

// Put reacts to database writes and implements dirty data uncaching. This is the
// post-processing step of a commit operation where the already persisted trie is
// removed from the dirty cache and moved into the clean cache. The reason behind
// the two-phase commit is to ensure ensure data availability while moving from
// memory to disk.
func (c *cleaner) Put(key []byte, rlp []byte) error {
	hash := types.BytesToHash(key)

	// If the node does not exist, we're done on this path
	node, ok := c.db.dirties[hash]
	if !ok {
		return nil
	}
	// Node still exists, remove it from the flush-list
	switch hash {
	case c.db.oldest:
		c.db.oldest = node.flushNext
		c.db.dirties[node.flushNext].flushPrev = types.Hash{}
	case c.db.newest:
		c.db.newest = node.flushPrev
		c.db.dirties[node.flushPrev].flushNext = types.Hash{}
	default:
		c.db.dirties[node.flushPrev].flushNext = node.flushNext
		c.db.dirties[node.flushNext].flushPrev = node.flushPrev
	}
	// Remove the node from the dirty cache
	delete(c.db.dirties, hash)
	c.db.dirtiesSize -= types.StorageSize(types.HashLength + int(node.size))
	if node.children != nil {
		c.db.dirtiesSize -= types.StorageSize(cachedNodeChildrenSize + len(node.children)*(types.HashLength+2))
	}
	// Move the flushed node into the clean cache to prevent insta-reloads
	if c.db.cleans != nil {
		c.db.cleans.Set(hash[:], rlp)
	}
	return nil
}

func (c *cleaner) Delete(key []byte) error {
	panic("not implemented")
}

// Size returns the current storage size of the memory cache in front of the
// persistent database layer.
func (db *Database) Size() (types.StorageSize, types.StorageSize) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	// db.dirtiesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted.
	var metadataSize = types.StorageSize((len(db.dirties) - 1) * cachedNodeSize)
	var metarootRefs = types.StorageSize(len(db.dirties[types.Hash{}].children) * (types.HashLength + 2))
	return db.dirtiesSize + db.childrenSize + metadataSize - metarootRefs, db.preimagesSize
}