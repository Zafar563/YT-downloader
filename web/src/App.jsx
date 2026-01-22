import { useState, useEffect, useRef, useMemo } from 'react';
import axios from 'axios';
import './index.css';

const API_BASE = 'http://localhost:8080/api';
const WS_URL = 'ws://localhost:8080/ws';

function App() {
  const [url, setUrl] = useState('');
  const [playlist, setPlaylist] = useState(null);
  const [loading, setLoading] = useState(false);
  const [selected, setSelected] = useState(new Set());
  const [format, setFormat] = useState('video'); // 'video' or 'mp3'
  const [progress, setProgress] = useState({}); // { videoId: { status, percent, message } }
  const [selectCount, setSelectCount] = useState(''); // New state for input number
  const wsRef = useRef(null);

  // Initialize WebSocket
  useEffect(() => {
    wsRef.current = new WebSocket(WS_URL);

    wsRef.current.onopen = () => console.log('WS Connected');
    wsRef.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // data = { video_id, status, percent, message }
      console.log("WS Data:", data);
      setProgress(prev => ({
        ...prev,
        [data.video_id]: data
      }));
    };

    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, []);

  const fetchPlaylist = async () => {
    if (!url) return;
    setLoading(true);
    try {
      const res = await axios.post(`${API_BASE}/playlist/info`, { url });
      setPlaylist(res.data);
      // Select all by default
      const allIds = new Set(res.data.entries.map(v => v.id || v.url || v.webpage_url));
      // Note: Video ID might be missing in flat-playlist, let's assume 'id' or fallback
      setSelected(allIds);
    } catch (err) {
      console.error(err);
      alert('Failed to fetch playlist: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  const toggleSelect = (id) => {
    const newSelected = new Set(selected);
    if (newSelected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelected(newSelected);
  };

  const selectFirstN = () => {
    const n = parseInt(selectCount, 10);
    if (isNaN(n) || n <= 0) {
      alert('Please enter a valid number');
      return;
    }

    if (!playlist || !playlist.entries) return;

    const newSelected = new Set();
    const limit = Math.min(n, playlist.entries.length);

    for (let i = 0; i < limit; i++) {
      const v = playlist.entries[i];
      newSelected.add(v.id || v.url || v.webpage_url);
    }

    setSelected(newSelected);
  };

  const startDownload = async () => {
    if (selected.size === 0) return;

    // Convert Set to Array of URLs
    // We need to find the full URL for each ID from the playlist entries
    const urlsToDownload = playlist.entries
      .filter(v => selected.has(v.id || v.url || v.webpage_url))
      .map(v => v.webpage_url || v.url); // yt-dlp flat playlist uses 'url' or 'webpage_url'

    try {
      await axios.post(`${API_BASE}/download`, { urls: urlsToDownload, format });
      alert('Download started!');
    } catch (err) {
      alert('Failed to start download');
    }
  };

  const formatDuration = (seconds) => {
    if (!seconds) return '0:00';
    const m = Math.floor(seconds / 60);
    const s = Math.floor(seconds % 60);
    return `${m}:${s.toString().padStart(2, '0')}`;
  };

  return (
    <div className="container">
      <header className="header">
        <h1>YT Downloader</h1>
        <p>Ultra-fast playlist downloader powered by Go & yt-dlp</p>
      </header>

      <div className="input-group">
        <input
          type="text"
          placeholder="Paste YouTube Playlist or Video URL"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button onClick={fetchPlaylist} disabled={loading}>
          {loading ? 'Fetching...' : 'Fetch'}
        </button>
      </div>

      {playlist && (
        <>
          <div className="playlist-info">
            <h2>{playlist.title || "Unknown Playlist"}</h2>
            <p>{playlist.entries.length} videos found</p>

            <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', alignItems: 'center', marginBottom: '1rem' }}>
              <button
                onClick={() => {
                  if (selected.size === playlist.entries.length) setSelected(new Set());
                  else setSelected(new Set(playlist.entries.map(v => v.id || v.url || v.webpage_url)));
                }}
                style={{ background: 'transparent', padding: '0', color: 'var(--accent)', textDecoration: 'underline' }}
              >
                {selected.size === playlist.entries.length ? 'Deselect All' : 'Select All'}
              </button>

              <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', marginLeft: 'auto' }}>
                <input
                  type="number"
                  min="1"
                  placeholder="Select N..."
                  value={selectCount}
                  onChange={(e) => setSelectCount(e.target.value)}
                  style={{ width: '100px', padding: '0.5rem', background: '#333', border: '1px solid #444', color: '#fff', borderRadius: '4px' }}
                />
                <button
                  onClick={selectFirstN}
                  style={{ padding: '0.5rem 1rem', fontSize: '0.9rem' }}
                >
                  Select First
                </button>
              </div>
            </div>

            <div style={{ marginTop: '0.5rem', display: 'flex', gap: '1rem', alignItems: 'center' }}>
              <span style={{ color: 'var(--text-sec)' }}>Format:</span>
              <button
                onClick={() => setFormat('video')}
                style={{
                  background: format === 'video' ? 'var(--primary)' : 'transparent',
                  color: format === 'video' ? '#000' : 'var(--text-main)',
                  border: '1px solid var(--primary)'
                }}
              >
                Default (Video)
              </button>
              <button
                onClick={() => setFormat('mp3')}
                style={{
                  background: format === 'mp3' ? 'var(--primary)' : 'transparent',
                  color: format === 'mp3' ? '#000' : 'var(--text-main)',
                  border: '1px solid var(--primary)'
                }}
              >
                MP3 Audio
              </button>
            </div>
          </div>

          <div className="video-grid">
            {playlist.entries.map((video, idx) => {
              // Fallback ID handling
              const vId = video.id || video.url || video.webpage_url || `vid-${idx}`;
              const vidProgress = progress[video.url] || progress[vId] || progress[video.webpage_url]; // Match by URL or ID used in backend

              return (
                <div key={vId} className="video-card">
                  <div className="thumbnail-wrapper">
                    {/* Checkbox Overlay */}
                    <div className="checkbox-wrapper">
                      <input
                        type="checkbox"
                        checked={selected.has(vId)}
                        onChange={() => toggleSelect(vId)}
                      />
                    </div>
                    {/* Thumbnail */}
                    {video.thumbnail ? (
                      <img src={video.thumbnail} alt={video.title} className="thumbnail" />
                    ) : (
                      <div className="thumbnail" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#333' }}>
                        No Preview
                      </div>
                    )}
                    <div className="duration">{formatDuration(video.duration)}</div>
                  </div>

                  <div className="card-content">
                    <h3 className="video-title" title={video.title}>{video.title || "Untitled Video"}</h3>

                    {vidProgress && (
                      <div className="progress-container">
                        <div className="progress-bar-bg">
                          <div
                            className="progress-bar-fill"
                            style={{ width: `${vidProgress.percent}%`, background: vidProgress.status === 'error' ? '#ff5555' : 'var(--accent)' }}
                          ></div>
                        </div>
                        <div className="status-text">
                          <span>{vidProgress.status === 'finished' ? 'Completed' : (vidProgress.status === 'error' ? 'Error' : `${vidProgress.percent.toFixed(1)}%`)}</span>
                          {vidProgress.message && <span style={{ color: '#ff5555' }}>{vidProgress.message}</span>}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              )
            })}
          </div>

          {selected.size > 0 && (
            <button className="download-fab" onClick={startDownload}>
              Download ({selected.size})
            </button>
          )}
        </>
      )}
    </div>
  );
}

export default App;
