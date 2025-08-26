// Gallery elements
const videoGrid = document.getElementById('videoGrid');
const loadingIndicator = document.getElementById('loadingIndicator');
const loadMoreBtn = document.getElementById('loadMoreBtn');
const searchInput = document.getElementById('searchInput');
const tagFilter = document.getElementById('tagFilter');
const sortBy = document.getElementById('sortBy');

// Modal elements
const videoModal = document.getElementById('videoModal');
const modalTitle = document.getElementById('modalTitle');
const modalVideo = document.getElementById('modalVideo');
const modalDescription = document.getElementById('modalDescription');
const modalTags = document.getElementById('modalTags');
const closeModal = document.getElementById('closeModal');

// State
let currentPage = 1;
let isLoading = false;
let hasMore = true;

// Load videos
async function loadVideos(page = 1, append = false) {
    if (isLoading || !hasMore) return;

    isLoading = true;
    loadingIndicator.classList.remove('hidden');
    loadMoreBtn.disabled = true;

    try {
        const search = searchInput.value;
        const tag = tagFilter.value;
        const sort = sortBy.value;

        const response = await fetch(`/api/videos?page=${page}&search=${search}&tag=${tag}&sort=${sort}`);
        if (!response.ok) {
            throw new Error('Failed to load videos');
        }

        const videos = await response.json();
        hasMore = videos.length === 10; // Assuming 10 videos per page

        if (!append) {
            videoGrid.innerHTML = '';
        }

        videos.forEach(video => {
            const videoCard = createVideoCard(video);
            videoGrid.appendChild(videoCard);
        });

        currentPage = page;
    } catch (err) {
        console.error('Error loading videos:', err);
    } finally {
        isLoading = false;
        loadingIndicator.classList.add('hidden');
        loadMoreBtn.disabled = false;
    }
}

// Create video card
function createVideoCard(video) {
    const card = document.createElement('div');
    card.className = 'group bg-white rounded-2xl shadow-lg hover:shadow-2xl transform hover:-translate-y-2 transition-all duration-300 overflow-hidden border border-gray-100';

    card.innerHTML = `
        <div class="relative aspect-w-16 aspect-h-9 cursor-pointer overflow-hidden">
            <img src="/api/thumbnails?id=${video.id}"
                 alt="${video.title}"
                 class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                 onerror="this.onerror=null; this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgdmlld0JveD0iMCAwIDMyMCAxODAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIzMjAiIGhlaWdodD0iMTgwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0xNjAgOTBDMTQzLjQzMSA5MCAxMzAgMTAzLjQzMSAxMzAgMTIwQzEzMCAxMzYuNTY5IDE0My40MzEgMTUwIDE2MCAxNTBDMTc2LjU2OSAxNTAgMTkwIDEzNi41NjkgMTkwIDEyMEMxOTAgMTAzLjQzMSAxNzYuNTY5IDkwIDE2MCA5MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPHBhdGggZD0iTTE2MCAxMzBDMTU1LjU4MiAxMzAgMTUyIDEyNi40MTggMTUyIDEyMkMxNTIgMTE3LjU4MiAxNTUuNTgyIDExNCAxNjAgMTE0QzE2NC40MTggMTE0IDE2OCAxMTcuNTgyIDE2OCAxMjJDMTY4IDEyNi40MTggMTY0LjQxOCAxMzAgMTYwIDEzMFoiIGZpbGw9IndoaXRlIi8+Cjwvc3ZnPgo='; this.alt='Video thumbnail not available';">

            <!-- Play button overlay -->
            <div class="absolute inset-0 bg-black/20 group-hover:bg-black/30 transition-colors duration-300 flex items-center justify-center">
                <div class="w-16 h-16 bg-white/90 group-hover:bg-white rounded-full flex items-center justify-center transform group-hover:scale-110 transition-all duration-300">
                    <svg class="w-8 h-8 text-indigo-600 ml-1" fill="currentColor" viewBox="0 0 24 24">
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
            <h3 class="font-bold text-gray-900 text-lg mb-2 line-clamp-2 group-hover:text-indigo-600 transition-colors duration-200">
                ${video.title}
            </h3>
            <p class="text-gray-600 text-sm mb-4 line-clamp-2 leading-relaxed">
                ${video.description}
            </p>

            <div class="flex flex-wrap gap-2 mb-4">
                ${video.tags.map(tag => `
                    <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800 group-hover:bg-indigo-200 transition-colors duration-200">
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
                <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                    </svg>
                    Video
                </span>
            </div>
        </div>
    `;

    // Add click handler to open modal
    card.querySelector('.aspect-w-16').addEventListener('click', () => {
        openVideoModal(video);
    });

    return card;
}

// Format duration from seconds to MM:SS
function formatDuration(seconds) {
    if (!seconds) return '';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
}

// Format date
function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 1) return 'Today';
    if (diffDays === 2) return 'Yesterday';
    if (diffDays <= 7) return `${diffDays - 1} days ago`;
    if (diffDays <= 30) return `${Math.ceil(diffDays / 7)} weeks ago`;
    if (diffDays <= 365) return `${Math.ceil(diffDays / 30)} months ago`;
    return date.toLocaleDateString();
}

// Open video modal
function openVideoModal(video) {
    modalTitle.textContent = video.title;
    modalVideo.src = `/api/videos/${video.id}`;
    modalDescription.textContent = video.description;

    modalTags.innerHTML = video.tags.map(tag => `
        <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
            ${tag}
        </span>
    `).join('');

    videoModal.classList.remove('hidden');

    // Use the new animation function if available
    if (window.openVideoModal) {
        window.openVideoModal(video);
    } else {
        // Fallback for older browsers
        modalVideo.play();
    }
}

// Close video modal
function closeVideoModal() {
    modalVideo.pause();
    modalVideo.src = '';

    // Use the new animation function if available
    if (window.closeVideoModal) {
        window.closeVideoModal();
    } else {
        // Fallback for older browsers
        videoModal.classList.add('hidden');
    }
}

// Event listeners
loadMoreBtn.addEventListener('click', () => {
    loadVideos(currentPage + 1, true);
});

searchInput.addEventListener('input', debounce(() => {
    currentPage = 1;
    loadVideos(1);
}, 300));

tagFilter.addEventListener('change', () => {
    currentPage = 1;
    loadVideos(1);
});

sortBy.addEventListener('change', () => {
    currentPage = 1;
    loadVideos(1);
});

closeModal.addEventListener('click', closeVideoModal);
videoModal.addEventListener('click', (e) => {
    if (e.target === videoModal) {
        closeVideoModal();
    }
});

// Utility function for debouncing
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Initialize
loadVideos(1);
