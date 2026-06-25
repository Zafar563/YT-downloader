import { useState } from 'react';
import axios from 'axios';
import './index.css';

const API_BASE = import.meta.env.VITE_API_URL || '/api';

// Icon Components for premium look
const SearchIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="11" cy="11" r="8"></circle>
    <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
  </svg>
);

const YoutubeIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" style={{ marginRight: '4px' }}>
    <path d="M23.498 6.163a3.003 3.003 0 0 0-2.11-2.107C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.388.511a3.002 3.002 0 0 0-2.11 2.107C0 8.053 0 12 0 12s0 3.947.502 5.837a3.003 3.003 0 0 0 2.11 2.107C4.495 20.455 12 20.455 12 20.455s7.505 0 9.388-.511a3.003 3.003 0 0 0 2.11-2.107C24 15.947 24 12 24 12s0-3.947-.502-5.837zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
  </svg>
);

const InstagramIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" style={{ marginRight: '4px' }}>
    <rect x="2" y="2" width="20" height="20" rx="5" ry="5"></rect>
    <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"></path>
    <line x1="17.5" y1="6.5" x2="17.51" y2="6.5"></line>
  </svg>
);

const DownloadIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
    <polyline points="7 10 12 15 17 10"></polyline>
    <line x1="12" y1="15" x2="12" y2="3"></line>
  </svg>
);

const ChevronIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="6 9 12 15 18 9"></polyline>
  </svg>
);

const CheckIcon = () => (
  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3.5" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="20 6 9 17 4 12"></polyline>
  </svg>
);

function App() {
  const [url, setUrl] = useState('');
  const [playlist, setPlaylist] = useState(null);
  const [loading, setLoading] = useState(false);
  const [selected, setSelected] = useState(new Set());
  const [format, setFormat] = useState('video'); // 'video' or 'mp3'
  const [quality, setQuality] = useState('best'); // 'best', '1080p', etc.
  const [selectCount, setSelectCount] = useState('');

  const fetchPlaylist = async () => {
    if (!url) return;
    setPlaylist(null); // Clear previous results on new fetch
    setLoading(true);
    try {
      const res = await axios.post(`${API_BASE}/playlist/info`, { url });
      setPlaylist(res.data);
      const allIds = new Set(res.data.entries.map(v => v.id || v.url || v.webpage_url));
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
    const selectedVideos = playlist.entries.filter(v =>
      selected.has(v.id || v.url || v.webpage_url)
    );
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

  const getPlatform = (videoUrl) => {
    if (!videoUrl) return '';
    const lower = videoUrl.toLowerCase();
    if (lower.includes('youtube.com') || lower.includes('youtu.be')) return 'youtube';
    if (lower.includes('instagram.com')) return 'instagram';
    return '';
  };

  return (
    <div className="container">
      <header className="hero">
        <div className="hero-badge">
          <span>✨</span> Free & High-Speed Media Downloader
        </div>
        <h1>YT & IG Downloader</h1>
        <p>Save your favorite videos and audio tracks from YouTube and Instagram in high quality instantly.</p>
      </header>

      <div className="search-card">
        <div className="input-group">
          <div className="input-wrapper">
            <span className="input-icon">
              <SearchIcon />
            </span>
            <input
              type="text"
              placeholder="Paste YouTube or Instagram link..."
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && fetchPlaylist()}
            />
          </div>
          <button className="btn-fetch" onClick={fetchPlaylist} disabled={loading}>
            {loading ? (
              <>
                <svg className="spinner" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3.5" strokeLinecap="round">
                  <circle cx="12" cy="12" r="10" strokeDasharray="32" strokeDashoffset="8"></circle>
                </svg>
                <span>Fetching...</span>
              </>
            ) : (
              'Fetch'
            )}
          </button>
        </div>
      </div>

      {loading && (
        <div className="loading-skeleton-container">
          <svg className="spinner" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="var(--primary)" strokeWidth="3" strokeLinecap="round">
            <circle cx="12" cy="12" r="10" strokeDasharray="32" strokeDashoffset="8"></circle>
          </svg>
          <div className="loading-text">Fetching media details from server...</div>
        </div>
      )}

      {playlist && (
        <>
          <div className="playlist-info">
            <div className="playlist-header-row">
              <div className="playlist-title-sec">
                <h2>{playlist.title || "Fetched Playlist"}</h2>
                <p>{playlist.entries.length} items found</p>
              </div>

              <div className="right-controls">
                <div className="config-group">
                  <span className="config-label">Format:</span>
                  <div className="toggle-group">
                    <button
                      className={`toggle-btn ${format === 'video' ? 'active' : ''}`}
                      onClick={() => setFormat('video')}
                    >
                      🎥 Video
                    </button>
                    <button
                      className={`toggle-btn ${format === 'mp3' ? 'active' : ''}`}
                      onClick={() => setFormat('mp3')}
                    >
                      🎵 MP3 Audio
                    </button>
                  </div>
                </div>

                {format === 'video' && (
                  <div className="config-group">
                    <span className="config-label">Quality:</span>
                    <div className="select-container">
                      <select
                        className="custom-select"
                        value={quality}
                        onChange={(e) => setQuality(e.target.value)}
                      >
                        <option value="best">Best Quality</option>
                        <option value="2160p">4K (2160p)</option>
                        <option value="1440p">2K (1440p)</option>
                        <option value="1080p">1080p</option>
                        <option value="720p">720p</option>
                        <option value="480p">480p</option>
                        <option value="360p">360p</option>
                      </select>
                      <span className="select-chevron">
                        <ChevronIcon />
                      </span>
                    </div>
                  </div>
                )}
              </div>
            </div>

            <div className="playlist-controls-row">
              <div className="left-controls">
                <button
                  className="text-link-btn"
                  onClick={() => {
                    if (selected.size === playlist.entries.length) setSelected(new Set());
                    else setSelected(new Set(playlist.entries.map(v => v.id || v.url || v.webpage_url)));
                  }}
                >
                  {selected.size === playlist.entries.length ? 'Deselect All' : 'Select All'}
                </button>
              </div>

              <div className="right-controls">
                <div className="select-n-group">
                  <input
                    type="number"
                    min="1"
                    placeholder="Count"
                    className="input-number-mini"
                    value={selectCount}
                    onChange={(e) => setSelectCount(e.target.value)}
                  />
                  <button className="btn-mini" onClick={selectFirstN}>
                    Select First N
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div className="video-grid">
            {playlist.entries.map((video, idx) => {
              const vId = video.id || video.url || video.webpage_url || `vid-${idx}`;
              const isSelected = selected.has(vId);
              const platform = getPlatform(video.webpage_url || video.url || url);

              return (
                <div
                  key={vId}
                  className={`video-card ${isSelected ? 'selected' : ''}`}
                  onClick={() => toggleSelect(vId)}
                >
                  <div className="thumbnail-wrapper">
                    {/* Custom Checkbox Overlay */}
                    <div className="checkbox-container">
                      <div className={`custom-checkbox ${isSelected ? 'checked' : ''}`}>
                        <CheckIcon />
                      </div>
                    </div>

                    {/* Platform Badge Overlay */}
                    {platform && (
                      <span className={`platform-badge ${platform}`}>
                        {platform === 'youtube' ? <YoutubeIcon /> : <InstagramIcon />}
                        {platform}
                      </span>
                    )}

                    {/* Thumbnail */}
                    {video.thumbnail ? (
                      <img src={video.thumbnail} alt={video.title} className="thumbnail" loading="lazy" />
                    ) : (
                      <div className="thumbnail" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'rgba(255,255,255,0.03)', color: 'var(--text-sec)', fontSize: '0.9rem' }}>
                        No Preview
                      </div>
                    )}
                    <div className="duration">{formatDuration(video.duration)}</div>
                  </div>

                  <div className="card-content">
                    <h3 className="video-title" title={video.title}>
                      {video.title || "Untitled Video"}
                    </h3>
                  </div>
                </div>
              );
            })}
          </div>

          {selected.size > 0 && (
            <button className="download-fab" onClick={(e) => { e.stopPropagation(); startDownload(); }}>
              <DownloadIcon />
              Download Selected ({selected.size})
            </button>
          )}
        </>
      )}
    </div>
  );
}

export default App;
