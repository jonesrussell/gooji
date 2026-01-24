import type { Video, VideoListResponse } from './types';

// Gallery elements
const videoGrid = document.getElementById('videoGrid') as HTMLElement | null;
const loadingIndicator = document.getElementById('loadingIndicator') as HTMLElement | null;
const loadMoreBtn = document.getElementById('loadMoreBtn') as HTMLButtonElement | null;
const searchInput = document.getElementById('searchInput') as HTMLInputElement | null;
const tagFilter = document.getElementById('tagFilter') as HTMLSelectElement | null;
const sortBy = document.getElementById('sortBy') as HTMLSelectElement | null;

// Modal elements
const videoModal = document.getElementById('videoModal') as HTMLElement | null;
const modalTitle = document.getElementById('modalTitle') as HTMLElement | null;
const modalVideo = document.getElementById('modalVideo') as HTMLVideoElement | null;
const modalDescription = document.getElementById('modalDescription') as HTMLElement | null;
const modalTags = document.getElementById('modalTags') as HTMLElement | null;
const closeModal = document.getElementById('closeModal') as HTMLElement | null;

// State
let currentPage = 1;
let isLoading = false;
let hasMore = true;

// Load videos
async function loadVideos(page = 1, append = false): Promise<void> {
  if (isLoading || !hasMore || !videoGrid || !loadingIndicator || !loadMoreBtn) return;

  isLoading = true;
  loadingIndicator.classList.remove('hidden');
  loadMoreBtn.disabled = true;

  try {
    const search = searchInput?.value || '';
    const tag = tagFilter?.value || '';
    const sort = sortBy?.value || '';

    const response = await fetch(`/api/videos?page=${page}&search=${encodeURIComponent(search)}&tag=${encodeURIComponent(tag)}&sort=${encodeURIComponent(sort)}`);
    if (!response.ok) {
      throw new Error('Failed to load videos');
    }

    const videos = await response.json() as VideoListResponse;
    hasMore = videos.length === 10; // Assuming 10 videos per page

    if (!append) {
      videoGrid.innerHTML = '';
    }

    videos.forEach((video: Video) => {
      const videoCard = createVideoCard(video);
      videoGrid.appendChild(videoCard);
    });

    currentPage = page;
  } catch (err) {
    console.error('Error loading videos:', err);
  } finally {
    isLoading = false;
    if (loadingIndicator) {
      loadingIndicator.classList.add('hidden');
    }
    if (loadMoreBtn) {
      loadMoreBtn.disabled = false;
    }
  }
}

// Create video card
function createVideoCard(video: Video): HTMLElement {
  const card = document.createElement('div');
  card.className = 'group bg-white rounded-2xl shadow-lg hover:shadow-2xl transform hover:-translate-y-2 transition-all duration-300 overflow-hidden border border-gray-100';
  card.setAttribute('data-video-id', video.id);

  const thumbnailUrl = `/api/thumbnails?id=${video.id}`;
  const fallbackThumbnail = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgdmlld0JveD0iMCAwIDMyMCAxODAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIzMjAiIGhlaWdodD0iMTgwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0xNjAgOTBDMTQzLjQzMSA5MCAxMzAgMTAzLjQzMSAxMzAgMTIwQzEzMCAxMzYuNTY5IDE0My40MzEgMTUwIDE2MCAxNTBDMTc2LjU2OSAxNTAgMTkwIDEzNi41NjkgMTkwIDEyMEMxOTAgMTAzLjQzMSAxNzYuNTY5IDkwIDE2MCA5MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPHBhdGggZD0iTTE2MCAxMzBDMTU1LjU4MiAxMzAgMTUyIDEyNi40MTggMTUyIDEyMkMxNTIgMTE3LjU4MiAxNTUuNTgyIDExNCAxNjAgMTE0QzE2NC40MTggMTE0IDE2OCAxMTcuNTgyIDE2OCAxMjJDMTY4IDEyNi40MTggMTY0LjQxOCAxMzAgMTYwIDEzMFoiIGZpbGw9IndoaXRlIi8+Cjwvc3ZnPgo=';

  card.innerHTML = `
    <div class="relative aspect-w-16 aspect-h-9 cursor-pointer overflow-hidden">
      <img src="${thumbnailUrl}"
           alt="${video.title}"
           class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
           onerror="this.onerror=null; this.src='${fallbackThumbnail}'; this.alt='Video thumbnail not available';">

      <!-- Play button overlay -->
      <div class="absolute inset-0 bg-black/20 group-hover:bg-black/30 transition-colors duration-300 flex items-center justify-center">
        <div class="w-16 h-16 bg-white/90 group-hover:bg-white rounded-full flex items-center justify-center transform group-hover:scale-110 transition-all duration-300">
          <svg class="w-8 h-8 text-yellow-600 ml-1" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </div>
      </div>

      <!-- Duration badge -->
      ${video.duration ? `
        <div class="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded-lg font-medium">
          ${formatDuration(video.duration)}
        </div>
      ` : ''}
    </div>

    <div class="p-6">
      <h3 class="font-bold text-gray-900 text-lg mb-2 line-clamp-2 group-hover:text-yellow-600 transition-colors duration-200">
        ${video.title}
      </h3>
      <p class="text-gray-600 text-sm mb-4 line-clamp-2 leading-relaxed">
        ${video.description}
      </p>

      <div class="flex flex-wrap gap-2 mb-4">
        ${video.tags.map(tag => `
          <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 group-hover:bg-yellow-200 transition-colors duration-200">
            ${tag}
          </span>
        `).join('')}
      </div>

      <div class="flex items-center justify-between text-xs text-gray-500">
        <span class="flex items-center">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          ${video.createdAt ? formatDate(video.createdAt) : 'Recently added'}
        </span>
        <div class="flex items-center space-x-2">
          <span class="flex items-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
            </svg>
            Video
          </span>
          <button onclick="deleteVideo('${video.id}', '${video.title.replace(/'/g, "\\'")}')"
                  class="flex items-center text-red-500 hover:text-red-700 hover:bg-red-50 p-1 rounded transition-colors duration-200"
                  title="Delete video">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>
  `;

  // Add click handler to open modal
  const thumbnailContainer = card.querySelector('.aspect-w-16');
  if (thumbnailContainer) {
    thumbnailContainer.addEventListener('click', () => {
      openVideoModal(video);
    });
  }

  return card;
}

// Format duration from seconds to MM:SS
function formatDuration(seconds: number): string {
  if (!seconds) return '';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

// Format date
function formatDate(dateString: string): string {
  if (!dateString) return '';
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = Math.abs(now.getTime() - date.getTime());
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 1) return 'Today';
  if (diffDays === 2) return 'Yesterday';
  if (diffDays <= 7) return `${diffDays - 1} days ago`;
  if (diffDays <= 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
  if (diffDays <= 365) return `${Math.ceil(diffDays / 30)} months ago`;
  return date.toLocaleDateString();
}

// Open video modal
function openVideoModal(video: Video): void {
  if (!modalTitle || !modalVideo || !modalDescription || !modalTags || !videoModal) return;

  modalTitle.textContent = video.title;
  modalVideo.src = `/api/videos/${video.id}`;
  modalDescription.textContent = video.description;

  modalTags.innerHTML = video.tags.map(tag => `
    <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800">
      ${tag}
    </span>
  `).join('');

  videoModal.classList.remove('hidden');

  // Use the new animation function if available
  if ((window as unknown as { openVideoModal?: (video: Video) => void }).openVideoModal) {
    (window as unknown as { openVideoModal: (video: Video) => void }).openVideoModal(video);
  } else {
    // Fallback for older browsers
    modalVideo.play();
  }
}

// Close video modal
function closeVideoModal(): void {
  if (!modalVideo || !videoModal) return;

  modalVideo.pause();
  modalVideo.src = '';

  // Use the new animation function if available
  if ((window as unknown as { closeVideoModal?: () => void }).closeVideoModal) {
    (window as unknown as { closeVideoModal: () => void }).closeVideoModal();
  } else {
    // Fallback for older browsers
    videoModal.classList.add('hidden');
  }
}

// Event listeners
if (loadMoreBtn) {
  loadMoreBtn.addEventListener('click', () => {
    loadVideos(currentPage + 1, true);
  });
}

if (searchInput) {
  searchInput.addEventListener('input', debounce(() => {
    currentPage = 1;
    loadVideos(1);
  }, 300));
}

if (tagFilter) {
  tagFilter.addEventListener('change', () => {
    currentPage = 1;
    loadVideos(1);
  });
}

if (sortBy) {
  sortBy.addEventListener('change', () => {
    currentPage = 1;
    loadVideos(1);
  });
}

if (closeModal) {
  closeModal.addEventListener('click', closeVideoModal);
}

if (videoModal) {
  videoModal.addEventListener('click', (e: MouseEvent) => {
    if (e.target === videoModal) {
      closeVideoModal();
    }
  });
}

// Utility function for debouncing
function debounce<T extends (...args: unknown[]) => void>(func: T, wait: number): (...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout> | null = null;
  return function executedFunction(...args: Parameters<T>): void {
    const later = (): void => {
      if (timeout) clearTimeout(timeout);
      func(...args);
    };
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// Delete video function
async function deleteVideo(videoId: string, videoTitle: string): Promise<void> {
  // Show confirmation dialog
  const confirmed = confirm(`Are you sure you want to delete "${videoTitle}"?\n\nThis action cannot be undone.`);

  if (!confirmed) {
    return;
  }

  try {
    const response = await fetch(`/api/videos/${videoId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to delete video');
    }

    const result = await response.json() as { message?: string };
    console.log('Video deleted:', result);

    // Remove the video card from the UI
    const videoCard = document.querySelector(`[data-video-id="${videoId}"]`);
    if (videoCard) {
      videoCard.remove();
    } else {
      // If we can't find the specific card, reload the videos
      loadVideos(1);
    }

    // Show success message
    showNotification('Video deleted successfully', 'success');
  } catch (error) {
    console.error('Error deleting video:', error);
    showNotification('Failed to delete video. Please try again.', 'error');
  }
}

// Show notification function
function showNotification(message: string, type: 'success' | 'error' | 'info' = 'info'): void {
  // Create notification element
  const notification = document.createElement('div');
  const bgColor = type === 'success' ? 'bg-green-500' : type === 'error' ? 'bg-red-500' : 'bg-blue-500';
  notification.className = `fixed top-4 right-4 z-50 px-6 py-4 rounded-lg shadow-lg transform transition-all duration-300 translate-x-full ${bgColor} text-white`;
  notification.textContent = message;

  // Add to page
  document.body.appendChild(notification);

  // Animate in
  setTimeout(() => {
    notification.classList.remove('translate-x-full');
  }, 100);

  // Remove after 3 seconds
  setTimeout(() => {
    notification.classList.add('translate-x-full');
    setTimeout(() => {
      if (notification.parentNode) {
        notification.parentNode.removeChild(notification);
      }
    }, 300);
  }, 3000);
}

// Make deleteVideo available globally for onclick handlers
(window as unknown as { deleteVideo: (videoId: string, videoTitle: string) => Promise<void> }).deleteVideo = deleteVideo;

// Initialize
loadVideos(1);
