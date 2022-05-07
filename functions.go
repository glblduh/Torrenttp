/* Contains all of the functions of the program */

package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

// Function for sending error message as JSON response
func errorRes(w http.ResponseWriter, error string, code int) {
	err := encodeRes(w, &jsonErrorRes{
		Error: error,
	})
	if err != nil {
		w.Write([]byte(error))
	}
}

// Compiles infohash, display name, and trackers to *torrent.TorrentSpec
func makeTorrentSpec(infohash string, displayname string, trackers []string) *torrent.TorrentSpec {
	spec := torrent.TorrentSpec{}
	spec.InfoHash = metainfo.NewHashFromHex(infohash)
	spec.DisplayName = displayname
	for _, tracker := range trackers {
		spec.Trackers = append(spec.Trackers, []string{tracker})
	}
	return &spec
}

// Decodes r.Body as JSON and automatically creates an error response if error
func decodeBody(w http.ResponseWriter, body io.Reader, v any) error {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		errorRes(w, "JSON Decoder error", http.StatusInternalServerError)
	}
	return err
}

// Encodes a struct as JSON response and automatically creates an error response if error
func encodeRes(w http.ResponseWriter, v any) error {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		errorRes(w, "JSON Encoder error", http.StatusInternalServerError)
	}
	return err
}

// Turns persistentSpec to *torrent.TorrentSpec
func persistSpecToTorrentSpec(spec persistentSpec) *torrent.TorrentSpec {
	return &torrent.TorrentSpec{
		Trackers:                 spec.Trackers,
		InfoHash:                 metainfo.NewHashFromHex(spec.InfoHash),
		DisplayName:              spec.DisplayName,
		Webseeds:                 spec.Webseeds,
		DhtNodes:                 spec.DhtNodes,
		PeerAddrs:                spec.PeerAddrs,
		Sources:                  spec.Sources,
		DisableInitialPieceCheck: spec.DisableInitialPieceCheck,
		DisallowDataUpload:       spec.DisallowDataUpload,
		DisallowDataDownload:     spec.DisallowDataDownload,
	}
}

// Creates a URL for the stream of file
func createFileLink(infohash string, filename string) string {
	return "/api/getfile?infohash=" + infohash + "&file=" + url.QueryEscape(filename)
}

// Get the file handle inside the torrent
func getTorrentFile(t *torrent.Torrent, displaypath string) (*torrent.File, error) {
	for _, f := range t.Files() {
		if f.DisplayPath() == displaypath {
			return f, nil
		}
	}
	return nil, errors.New("File not found")
}

// Create config for BitTorrent client with confs from args
func newBtCliConfs(dir string, noup bool) *torrent.ClientConfig {
	opts := torrent.NewDefaultClientConfig()
	opts.DataDir = dir
	opts.NoUpload = noup
	return opts
}
