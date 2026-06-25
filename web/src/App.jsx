import { useState, useEffect } from 'react';
import axios from 'axios';
import './index.css';

const API_BASE = import.meta.env.VITE_API_URL || '/api';

const AVATAR_COLORS = ['#6366f1', '#ec4899', '#f43f5e', '#10b981', '#f59e0b', '#06b6d4'];

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

const UserIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
    <circle cx="12" cy="7" r="4"></circle>
  </svg>
);

const CloseIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <line x1="18" y1="6" x2="6" y2="18"></line>
    <line x1="6" y1="6" x2="18" y2="18"></line>
  </svg>
);

const TrashIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="3 6 5 6 21 6"></polyline>
    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
  </svg>
);

const SettingsIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="3"></circle>
    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
  </svg>
);

const ClockIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="10"></circle>
    <polyline points="12 6 12 12 16 14"></polyline>
  </svg>
);

const LogoutIcon = () => (
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
    <polyline points="16 17 21 12 16 7"></polyline>
    <line x1="21" y1="12" x2="9" y2="12"></line>
  </svg>
);

function App() {
  // Auth Token
  const [token, setToken] = useState(() => localStorage.getItem('yt_ig_auth_token') || '');

  // User Profile
  const [profile, setProfile] = useState(() => {
    const saved = localStorage.getItem('yt_ig_user_profile');
    if (saved) {
      try {
        return JSON.parse(saved);
      } catch (e) {}
    }
    return {
      name: 'Guest User',
      email: '',
      defaultFormat: 'video',
      defaultQuality: 'best',
      avatarColor: '#6366f1'
    };
  });

  // Client History logs
  const [history, setHistory] = useState(() => {
    const saved = localStorage.getItem('yt_ig_download_history');
    if (saved) {
      try {
        return JSON.parse(saved);
      } catch (e) {}
    }
    return [];
  });

  // App core states
  const [url, setUrl] = useState('');
  const [playlist, setPlaylist] = useState(null);
  const [loading, setLoading] = useState(false);
  const [selected, setSelected] = useState(new Set());
  const [format, setFormat] = useState(profile.defaultFormat || 'video');
  const [quality, setQuality] = useState(profile.defaultQuality || 'best');
  const [selectCount, setSelectCount] = useState('');
  const [isProfileOpen, setIsProfileOpen] = useState(false);

  // Drawer Auth Screen States
  const [authTab, setAuthTab] = useState('login'); // 'login' or 'register'
  const [authUsername, setAuthUsername] = useState('');
  const [authEmail, setAuthEmail] = useState('');
  const [authPassword, setAuthPassword] = useState('');
  const [authError, setAuthError] = useState('');
  const [authSuccess, setAuthSuccess] = useState('');
  const [authLoading, setAuthLoading] = useState(false);

  // Fetch Profile settings & History from Server if Token exists
  useEffect(() => {
    if (token) {
      fetchServerProfile();
    }
  }, [token]);

  const fetchServerProfile = async () => {
    try {
      const res = await axios.get(`${API_BASE}/auth/profile`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      const serverUser = res.data.user;
      const serverHistory = res.data.history;

      // Update Profile states
      const updatedProfile = {
        name: serverUser.username,
        email: serverUser.email,
        defaultFormat: serverUser.default_format || 'video',
        defaultQuality: serverUser.default_quality || 'best',
        avatarColor: serverUser.avatar_color || '#6366f1'
      };
      setProfile(updatedProfile);
      localStorage.setItem('yt_ig_user_profile', JSON.stringify(updatedProfile));
      setFormat(updatedProfile.defaultFormat);
      setQuality(updatedProfile.defaultQuality);

      // Map Server History to matching Local Format
      if (serverHistory && Array.isArray(serverHistory)) {
        const mappedHistory = serverHistory.map(item => ({
          id: item.id.toString(),
          title: item.url, // Fallback title
          url: item.url,
          platform: item.platform,
          format: item.format,
          date: new Date(item.downloaded_at).toLocaleDateString(undefined, {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
          })
        }));
        setHistory(mappedHistory);
        localStorage.setItem('yt_ig_download_history', JSON.stringify(mappedHistory));
      }
    } catch (err) {
      console.error('Failed to sync profile with server:', err);
      if (err.response?.status === 401) {
        // Token expired/invalid
        handleLogout();
      }
    }
  };

  const updateProfile = async (updated) => {
    setProfile(updated);
    localStorage.setItem('yt_ig_user_profile', JSON.stringify(updated));
    setFormat(updated.defaultFormat);
    setQuality(updated.defaultQuality);

    if (token) {
      try {
        await axios.put(
          `${API_BASE}/auth/profile`,
          {
            username: updated.name,
            email: updated.email,
            avatar_color: updated.avatarColor,
            default_format: updated.defaultFormat,
            default_quality: updated.defaultQuality
          },
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );
      } catch (err) {
        console.error('Failed to save profile preferences to database:', err);
      }
    }
  };

  const addToHistory = async (title, videoUrl) => {
    const platform = getPlatform(videoUrl);
    const dateStr = new Date().toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });

    const newEntry = {
      id: Date.now().toString(),
      title: title || 'Media File',
      url: videoUrl,
      platform: platform || 'unknown',
      format: format,
      date: dateStr
    };

    const updated = [newEntry, ...history].slice(0, 50);
    setHistory(updated);
    localStorage.setItem('yt_ig_download_history', JSON.stringify(updated));

    if (token) {
      try {
        await axios.post(
          `${API_BASE}/auth/download`,
          {
            platform: platform || 'unknown',
            url: videoUrl,
            format: format
          },
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );
      } catch (err) {
        console.error('Failed to log download to server:', err);
      }
    }
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    setAuthError('');
    setAuthSuccess('');
    setAuthLoading(true);

    try {
      const res = await axios.post(`${API_BASE}/auth/register`, {
        username: authUsername,
        email: authEmail,
        password: authPassword
      });

      const data = res.data;
      setToken(data.token);
      localStorage.setItem('yt_ig_auth_token', data.token);

      const serverProfile = {
        name: data.user.username,
        email: data.user.email,
        defaultFormat: data.user.default_format || 'video',
        defaultQuality: data.user.default_quality || 'best',
        avatarColor: data.user.avatar_color || '#6366f1'
      };
      setProfile(serverProfile);
      localStorage.setItem('yt_ig_user_profile', JSON.stringify(serverProfile));

      setAuthSuccess('Account registered successfully!');
      setAuthUsername('');
      setAuthEmail('');
      setAuthPassword('');
    } catch (err) {
      console.error(err);
      setAuthError(err.response?.data?.error || 'Registration failed. Please try again.');
    } finally {
      setAuthLoading(false);
    }
  };

  const handleLogin = async (e) => {
    e.preventDefault();
    setAuthError('');
    setAuthSuccess('');
    setAuthLoading(true);

    try {
      const res = await axios.post(`${API_BASE}/auth/login`, {
        username_or_email: authUsername,
        password: authPassword
      });

      const data = res.data;
      setToken(data.token);
      localStorage.setItem('yt_ig_auth_token', data.token);

      const serverProfile = {
        name: data.user.username,
        email: data.user.email,
        defaultFormat: data.user.default_format || 'video',
        defaultQuality: data.user.default_quality || 'best',
        avatarColor: data.user.avatar_color || '#6366f1'
      };
      setProfile(serverProfile);
      localStorage.setItem('yt_ig_user_profile', JSON.stringify(serverProfile));

      setAuthSuccess('Logged in successfully!');
      setAuthUsername('');
      setAuthPassword('');
    } catch (err) {
      console.error(err);
      setAuthError(err.response?.data?.error || 'Login failed. Invalid username or password.');
    } finally {
      setAuthLoading(false);
    }
  };

  const handleLogout = () => {
    setToken('');
    localStorage.removeItem('yt_ig_auth_token');
    
    const guestProfile = {
      name: 'Guest User',
      email: '',
      defaultFormat: 'video',
      defaultQuality: 'best',
      avatarColor: '#6366f1'
    };
    setProfile(guestProfile);
    localStorage.setItem('yt_ig_user_profile', JSON.stringify(guestProfile));
    setHistory([]);
    localStorage.removeItem('yt_ig_download_history');
    setAuthSuccess('');
    setAuthError('');
  };

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

      // Save to history log
      addToHistory(video.title || 'video', videoUrl);

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
        <div className="header-row">
          <div style={{ width: '42px' }}></div> {/* Centering spacer */}
          <div className="hero-badge">
            <span>✨</span> Free & High-Speed Media Downloader
          </div>
          <button className="profile-trigger-btn" onClick={() => setIsProfileOpen(true)} title="View User Profile">
            <div className="avatar-circle" style={{ backgroundColor: profile.avatarColor }}>
              {profile.name ? profile.name.charAt(0).toUpperCase() : 'G'}
            </div>
            <span className="profile-trigger-name">{profile.name}</span>
          </button>
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

      {/* Sliding Glassmorphic Profile Drawer */}
      {isProfileOpen && (
        <div className="drawer-overlay" onClick={() => setIsProfileOpen(false)}>
          <div className="drawer" onClick={(e) => e.stopPropagation()}>
            <div className="drawer-header">
              <h2>{token ? 'My Profile' : 'Authentication'}</h2>
              <button className="btn-close-drawer" onClick={() => setIsProfileOpen(false)} title="Close Panel">
                <CloseIcon />
              </button>
            </div>

            <div className="drawer-content">
              {/* Form alerts */}
              {authError && <div className="auth-error-msg">{authError}</div>}
              {authSuccess && <div className="auth-success-msg">{authSuccess}</div>}

              {/* Logged Out view: show Register/Login forms */}
              {!token ? (
                <div className="drawer-section">
                  <div className="auth-tabs">
                    <button
                      className={`auth-tab-btn ${authTab === 'login' ? 'active' : ''}`}
                      onClick={() => {
                        setAuthTab('login');
                        setAuthError('');
                        setAuthSuccess('');
                      }}
                    >
                      Login
                    </button>
                    <button
                      className={`auth-tab-btn ${authTab === 'register' ? 'active' : ''}`}
                      onClick={() => {
                        setAuthTab('register');
                        setAuthError('');
                        setAuthSuccess('');
                      }}
                    >
                      Register
                    </button>
                  </div>

                  {authTab === 'login' ? (
                    <form onSubmit={handleLogin} className="form-grid">
                      <div className="form-group">
                        <label htmlFor="login-username">Username or Email</label>
                        <input
                          id="login-username"
                          type="text"
                          value={authUsername}
                          onChange={(e) => setAuthUsername(e.target.value)}
                          placeholder="Enter your username or email..."
                          required
                        />
                      </div>

                      <div className="form-group">
                        <label htmlFor="login-pass">Password</label>
                        <input
                          id="login-pass"
                          type="password"
                          value={authPassword}
                          onChange={(e) => setAuthPassword(e.target.value)}
                          placeholder="••••••••"
                          required
                        />
                      </div>

                      <button type="submit" className="btn-auth-submit" disabled={authLoading}>
                        {authLoading ? 'Logging in...' : 'Sign In'}
                      </button>

                      <p className="auth-switch-text">
                        Don't have an account?
                        <span className="auth-switch-link" onClick={() => setAuthTab('register')}>
                          Sign Up
                        </span>
                      </p>
                    </form>
                  ) : (
                    <form onSubmit={handleRegister} className="form-grid">
                      <div className="form-group">
                        <label htmlFor="reg-username">Username</label>
                        <input
                          id="reg-username"
                          type="text"
                          value={authUsername}
                          onChange={(e) => setAuthUsername(e.target.value)}
                          placeholder="Pick a username..."
                          required
                        />
                      </div>

                      <div className="form-group">
                        <label htmlFor="reg-email">Email Address</label>
                        <input
                          id="reg-email"
                          type="email"
                          value={authEmail}
                          onChange={(e) => setAuthEmail(e.target.value)}
                          placeholder="Enter your email address..."
                          required
                        />
                      </div>

                      <div className="form-group">
                        <label htmlFor="reg-pass">Password</label>
                        <input
                          id="reg-pass"
                          type="password"
                          value={authPassword}
                          onChange={(e) => setAuthPassword(e.target.value)}
                          placeholder="At least 6 characters..."
                          required
                        />
                      </div>

                      <button type="submit" className="btn-auth-submit" disabled={authLoading}>
                        {authLoading ? 'Registering...' : 'Create Account'}
                      </button>

                      <p className="auth-switch-text">
                        Already have an account?
                        <span className="auth-switch-link" onClick={() => setAuthTab('login')}>
                          Sign In
                        </span>
                      </p>
                    </form>
                  )}
                </div>
              ) : (
                /* Logged In View */
                <>
                  {/* Personal Information */}
                  <div className="drawer-section">
                    <div className="drawer-section-title">
                      <UserIcon /> Personal Information
                    </div>

                    <div className="form-grid">
                      <div className="form-group" style={{ alignItems: 'center', marginBottom: '1rem' }}>
                        <div className="avatar-circle" style={{ backgroundColor: profile.avatarColor, width: '64px', height: '64px', fontSize: '1.75rem', border: '3px solid rgba(255,255,255,0.25)' }}>
                          {profile.name ? profile.name.charAt(0).toUpperCase() : 'G'}
                        </div>
                        <div className="avatar-presets-grid" style={{ marginTop: '0.75rem' }}>
                          {AVATAR_COLORS.map(color => (
                            <div
                              key={color}
                              className={`avatar-preset-option ${profile.avatarColor === color ? 'selected' : ''}`}
                              style={{ backgroundColor: color }}
                              onClick={() => updateProfile({ ...profile, avatarColor: color })}
                            />
                          ))}
                        </div>
                      </div>

                      <div className="form-group">
                        <label htmlFor="profile-name">Nickname</label>
                        <input
                          id="profile-name"
                          type="text"
                          value={profile.name}
                          onChange={(e) => updateProfile({ ...profile, name: e.target.value })}
                          placeholder="Enter your nickname..."
                        />
                      </div>

                      <div className="form-group">
                        <label htmlFor="profile-email">Email Address</label>
                        <input
                          id="profile-email"
                          type="email"
                          value={profile.email}
                          onChange={(e) => updateProfile({ ...profile, email: e.target.value })}
                          placeholder="yourname@example.com"
                          disabled // email can be readonly for security
                        />
                      </div>
                    </div>
                  </div>

                  {/* Preferences Settings */}
                  <div className="drawer-section">
                    <div className="drawer-section-title">
                      <SettingsIcon /> Default Preferences
                    </div>

                    <div className="form-grid">
                      <div className="form-group">
                        <label htmlFor="pref-format">Preferred Format</label>
                        <select
                          id="pref-format"
                          value={profile.defaultFormat}
                          onChange={(e) => updateProfile({ ...profile, defaultFormat: e.target.value })}
                        >
                          <option value="video">Default (Video)</option>
                          <option value="mp3">MP3 Audio</option>
                        </select>
                      </div>

                      <div className="form-group">
                        <label htmlFor="pref-quality">Preferred Quality</label>
                        <select
                          id="pref-quality"
                          value={profile.defaultQuality}
                          onChange={(e) => updateProfile({ ...profile, defaultQuality: e.target.value })}
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
                    </div>
                  </div>

                  {/* Local Download History */}
                  <div className="drawer-section" style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                    <div className="drawer-section-title" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <span style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <ClockIcon /> Download History (Synced)
                      </span>
                      {history.length > 0 && (
                        <span style={{ fontSize: '0.8rem', color: 'var(--text-sec)' }}>
                          Last {history.length} downloads
                        </span>
                      )}
                    </div>

                    {history.length === 0 ? (
                      <div className="history-empty">
                        <span>📂</span>
                        No download logs recorded yet.
                      </div>
                    ) : (
                      <div className="history-list" style={{ overflowY: 'auto', maxHeight: '250px' }}>
                        {history.map(item => (
                          <div key={item.id} className="history-item">
                            <div className="history-details">
                              <a href={item.url} target="_blank" rel="noopener noreferrer" className="history-title" title={item.title}>
                                {item.title}
                              </a>
                              <div className="history-meta">
                                <span className={`history-platform-tag ${item.platform}`}>{item.platform}</span>
                                <span className="history-format">{item.format.toUpperCase()}</span>
                                <span className="history-date">{item.date}</span>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Log Out button */}
                  <button className="btn-logout" onClick={handleLogout}>
                    <LogoutIcon />
                    Log Out
                  </button>
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
