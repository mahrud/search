package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"io/ioutil"
	"math"
	"os"
	"path"
	"search/index"
	"search/searcher"
	"search/util"
	"strconv"
)

// Server contains all the necessary information for a running server.
type Server struct {
	mountPoint string   // Mount point of the server
	lenMS      int      // Length of the master secret in bytes
	keyHalves  [][]byte // The server-side keyhalves
	salts      [][]byte // The salts for deriving the keys for the PRFs
	numFiles   int      // The number of files currently stored in the server.  This is used to determine the next docID.
}

// CreateServer initializes a server with `numClients` clients with a master
// secret of length `lenMS`, and generate salts with length `lenSalt`.  The
// number of salts is given by `r = -log2(fpRate)`, where `fpRate` is the
// desired false positive rate of the system.  `mountPoint` determines where the
// server files will be stored.
func CreateServer(numClients, lenMS, lenSalt int, mountPoint string, fpRate float64) *Server {
	s := new(Server)
	masterSecret := make([]byte, lenMS)
	rand.Read(masterSecret)
	s.keyHalves = make([][]byte, numClients)
	s.lenMS = lenMS
	for i := 0; i < numClients; i++ {
		h := sha256.New()
		h.Write([]byte(strconv.Itoa(i)))
		cksum := h.Sum(nil)
		s.keyHalves[i] = util.XorBytes(masterSecret, cksum, lenMS)
	}
	r := int(math.Ceil(-math.Log2(fpRate)))
	s.salts = util.GenerateSalts(r, lenSalt)
	s.numFiles = 0
	s.mountPoint = mountPoint
	s.writeToFile()
	return s
}

// LoadServer initializes a Server by reading the metadata stored at
// `mountPoint` and restoring the server status.
func LoadServer(mountPoint string) *Server {
	input, err := os.Open(path.Join(mountPoint, "serverMD"))
	if err != nil {
		panic("Server metadata not found")
	}
	dec := gob.NewDecoder(input)

	s := new(Server)
	dec.Decode(&s.mountPoint)
	dec.Decode(&s.numFiles)
	dec.Decode(&s.salts)
	dec.Decode(&s.keyHalves)
	dec.Decode(&s.lenMS)

	input.Close()

	return s
}

// writeToFile serializes the server status and writes the metadata to a file in
// the server mount point, which can be later loaded by `LoadServer`.
func (s *Server) writeToFile() {
	file, _ := os.Create(path.Join(s.mountPoint, "serverMD"))
	enc := gob.NewEncoder(file)
	enc.Encode(s.mountPoint)
	enc.Encode(s.numFiles)
	enc.Encode(s.salts)
	enc.Encode(s.keyHalves)
	enc.Encode(s.lenMS)

	file.Close()
}

// AddFile adds a file with `content` to the server with the document ID equal
// to the number of files currently in the server and updates the count.
// Returns the document ID.
func (s *Server) AddFile(content []byte) int {
	output, _ := os.Create(path.Join(s.mountPoint, strconv.Itoa(s.numFiles)))
	output.Write(content)
	s.numFiles++
	output.Close()
	s.writeToFile()
	return s.numFiles - 1
}

// GetFile returns the content of the document with `docID`.  Behavior is
// undefined if the docID is invalid (out of range).
func (s *Server) GetFile(docID int) []byte {
	content, _ := ioutil.ReadFile(path.Join(s.mountPoint, strconv.Itoa(docID)))
	return content
}

// WriteIndex writes a SecureIndex to the disk of the server.
func (s *Server) WriteIndex(si index.SecureIndex) {
	output := si.Marshal()
	file, _ := os.Create(path.Join(s.mountPoint, strconv.Itoa(si.DocID)+".index"))
	file.Write(output)
	file.Close()
}

// readIndex loads an index from the disk.
func (s *Server) readIndex(docID int) index.SecureIndex {
	input, _ := ioutil.ReadFile(path.Join(s.mountPoint, strconv.Itoa(docID)+".index"))
	si := index.Unmarshal(input)
	return si
}

// SearchWord searchers the server for a word with `trapdoors`.  Returns a list
// of document ids of files possibly containing the word in increasing order.
func (s *Server) SearchWord(trapdoors [][]byte) []int {
	var result []int
	for i := 0; i < s.numFiles; i++ {
		if searcher.SearchSecureIndex(s.readIndex(i), trapdoors) {
			result = append(result, i)
		}
	}
	return result
}

// WriteLookupTable writes `content` to the file "lookupTable".
func (s *Server) WriteLookupTable(content []byte) {
	file, _ := os.Create(path.Join(s.mountPoint, "lookupTable"))
	file.Write(content)
	file.Close()
}

// ReadLookupTable reads the content in the file "lookupTable" and returns it in
// a byte slice.
func (s *Server) ReadLookupTable() []byte {
	content, _ := ioutil.ReadFile(path.Join(s.mountPoint, "lookupTable"))
	return content
}

// GetKeyHalf returns the server-side key half for client with `clientNum`.
// Behavior is undefined if `clientNum` is invalid (out of range).
func (s *Server) GetKeyHalf(clientNum int) []byte {
	return s.keyHalves[clientNum]
}

// GetSalts returns the salts to the client.
func (s *Server) GetSalts() [][]byte {
	return s.salts
}
