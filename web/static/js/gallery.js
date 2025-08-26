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
    card.className = 'bg-white rounded-lg shadow-md overflow-hidden';
    card.innerHTML = `
        <div class="aspect-w-16 aspect-h-9 cursor-pointer">
            <img src="/api/thumbnails?id=${video.id}"
                 alt="${video.title}"
                 class="w-full h-full object-cover"
                 onerror="this.onerror=null; this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIwIiBoZWlnaHQ9IjE4MCIgdmlld0JveD0iMCAwIDMyMCAxODAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxyZWN0IHdpZHRoPSIzMjAiIGhlaWdodD0iMTgwIiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0xNjAgOTBDMTQzLjQzMSA5MCAxMzAgMTAzLjQzMSAxMzAgMTIwQzEzMCAxMzYuNTY5IDE0My40MzEgMTUwIDE2MCAxNTBDMTc2LjU2OSAxNTAgMTkwIDEzNi41NjkgMTkwIDEyMEMxOTAgMTAzLjQzMSAxNzYuNTY5IDkwIDE2MCA5MFoiIGZpbGw9IiM5Q0EzQUYiLz4KPHBhdGggZD0iTTE2MCAxMzBDMTU1LjU4MiAxMzAgMTUyIDEyNi40MTggMTUyIDEyMkMxNTIgMTE3LjU4MiAxNTUuNTgyIDExNCAxNjAgMTE0QzE2NC40MTggMTE0IDE2OCAxMTcuNTgyIDE2OCAxMjJDMTY4IDEyNi40MTggMTY0LjQxOCAxMzAgMTYwIDEzMFoiIGZpbGw9IndoaXRlIi8+Cjwvc3ZnPgo='; this.alt='Video thumbnail not available';">
        </div>
        <div class="p-4">
            <h3 class="font-semibold text-gray-800">${video.title}</h3>
            <p class="text-sm text-gray-600 mt-1">${video.description}</p>
            <div class="mt-2 flex flex-wrap gap-2">
                ${video.tags.map(tag => `
                    <span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">${tag}</span>
                `).join('')}
            </div>
        </div>
    `;

    // Add click handler to open modal
    card.querySelector('.aspect-w-16').addEventListener('click', () => {
        openVideoModal(video);
    });

    return card;
}

// Open video modal
function openVideoModal(video) {
    modalTitle.textContent = video.title;
    modalVideo.src = `/api/videos?id=${video.id}`;
    modalDescription.textContent = video.description;

    modalTags.innerHTML = video.tags.map(tag => `
        <span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">${tag}</span>
    `).join('');

    videoModal.classList.remove('hidden');
    modalVideo.play();
}

// Close video modal
function closeVideoModal() {
    videoModal.classList.add('hidden');
    modalVideo.pause();
    modalVideo.src = '';
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
