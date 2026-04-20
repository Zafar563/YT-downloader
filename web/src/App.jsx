import { useState, useEffect, useRef, useMemo } from 'react';
import axios from 'axios';
import './index.css';

const API_BASE = 'http://localhost:8080/api';

function App() {
  const [url, setUrl] = useState('');
  const [playlist, setPlaylist] = useState(null);
  const [loading, setLoading] = useState(false);
  const [selected, setSelected] = useState(new Set());
  const [format, setFormat] = useState('video'); // 'video' or 'mp3'
  const [quality, setQuality] = useState('best'); // 'best', '1080p', '720p', etc.
  const [selectCount, setSelectCount] = useState(''); // New state for input number


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

    // Convert Set to Array of needed data
    const selectedVideos = playlist.entries.filter(v =>
      selected.has(v.id || v.url || v.webpage_url)
    );

    // Trigger download for each selected video
    // Use hidden anchor tag to avoid opening new windows
    selectedVideos.forEach(video => {
      const videoUrl = video.webpage_url || video.url;
      const title = encodeURIComponent(video.title || 'video');
      const downloadLink = `${API_BASE}/stream?url=${encodeURIComponent(videoUrl)}&title=${title}&format=${format}&quality=${quality}`;

      const link = document.createElement('a');
      link.href = downloadLink;
      link.style.display = 'none';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    });
  };

  const formatDuration = (seconds) => {
    if (!seconds) return '0:00';
    const m = Math.floor(seconds / 60);
    const s = Math.floor(seconds % 60);
    return `${m}:${s.toString().padStart(2, '0')}`;
  };

  return (
    <div className="container">
      <header className="header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>YT Downloader</h1>
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

            {format === 'video' && (
              <div style={{ marginTop: '1rem', display: 'flex', gap: '1rem', alignItems: 'center' }}>
                <span style={{ color: 'var(--text-sec)' }}>Quality:</span>
                <select
                  value={quality}
                  onChange={(e) => setQuality(e.target.value)}
                  style={{
                    padding: '0.5rem',
                    background: '#333',
                    border: '1px solid var(--primary)',
                    color: '#fff',
                    borderRadius: '4px',
                    cursor: 'pointer'
                  }}
                >
                  <option value="best">Best Quality</option>
                  <option value="2160p">4K (2160p)</option>
                  <option value="1440p">2K (1440p)</option>
                  <option value="1080p">1080p</option>
                  <option value="720p">720p</option>
                  <option value="480p">480p</option>
                  <option value="360p">360p</option>
                </select>
              </div>
            )}
          </div>

          <div className="video-grid">
            {playlist.entries.map((video, idx) => {
              // Fallback ID handling
              const vId = video.id || video.url || video.webpage_url || `vid-${idx}`;
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
