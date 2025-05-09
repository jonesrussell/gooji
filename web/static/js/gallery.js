// Gallery functionality
class VideoGallery {
    constructor() {
        this.videos = [];
        this.page = 1;
        this.loading = false;
        this.filters = {
            search: '',
            category: '',
            sort: 'newest'
        };

        // DOM elements
        this.videoGrid = document.getElementById('video-grid');
        this.loadingIndicator = document.getElementById('loading');
        this.loadMoreBtn = document.getElementById('load-more');
        this.searchInput = document.getElementById('search');
        this.categorySelect = document.getElementById('category');
        this.sortSelect = document.getElementById('sort');
        this.modal = document.getElementById('video-modal');
        this.modalVideo = document.getElementById('modal-video');
        this.modalTitle = document.getElementById('modal-title');
        this.modalDescription = document.getElementById('modal-description');
        this.closeModalBtn = document.getElementById('close-modal');

        // Bind event handlers
        this.loadMoreBtn.addEventListener('click', () => this.loadMore());
        this.searchInput.addEventListener('input', () => this.handleSearch());
        this.categorySelect.addEventListener('change', () => this.handleFilter());
        this.sortSelect.addEventListener('change', () => this.handleFilter());
        this.closeModalBtn.addEventListener('click', () => this.closeModal());

        // Initialize
        this.loadVideos();
    }

    async loadVideos() {
        if (this.loading) return;

        this.loading = true;
        this.loadingIndicator.classList.remove('hidden');
        this.loadMoreBtn.classList.add('hidden');

        try {
            const response = await fetch(`/api/videos?page=${this.page}&${new URLSearchParams(this.filters)}`);
            if (!response.ok) throw new Error('Failed to load videos');

            const videos = await response.json();
            this.videos = this.page === 1 ? videos : [...this.videos, ...videos];
            this.renderVideos();

            this.loadMoreBtn.classList.remove('hidden');
        } catch (error) {
            console.error('Error loading videos:', error);
            // Show error message to user
        } finally {
            this.loading = false;
            this.loadingIndicator.classList.add('hidden');
        }
    }

    renderVideos() {
        const start = (this.page - 1) * 12;
        const end = start + 12;
        const videosToShow = this.videos.slice(start, end);

        videosToShow.forEach(video => {
            const card = this.createVideoCard(video);
            this.videoGrid.appendChild(card);
        });
    }

    createVideoCard(video) {
        const card = document.createElement('div');
        card.className = 'bg-white rounded-lg shadow-lg overflow-hidden';
        card.innerHTML = `
            <div class="aspect-w-16 aspect-h-9">
                <img src="${video.thumbnail}" alt="${video.title}" class="w-full h-full object-cover">
                <div class="absolute inset-0 bg-black bg-opacity-0 hover:bg-opacity-20 transition-opacity flex items-center justify-center">
                    <button class="play-button bg-red-600 text-white p-3 rounded-full opacity-0 hover:opacity-100 transition-opacity">
                        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                    </button>
                </div>
            </div>
            <div class="p-4">
                <h3 class="text-lg font-semibold mb-2">${video.title}</h3>
                <p class="text-gray-600 text-sm mb-2">${video.description}</p>
                <div class="flex justify-between items-center">
                    <span class="text-sm text-gray-500">${new Date(video.created_at).toLocaleDateString()}</span>
                    <span class="text-sm text-gray-500">${video.duration}s</span>
                </div>
            </div>
        `;

        // Add click handler
        card.querySelector('.play-button').addEventListener('click', () => this.openModal(video));

        return card;
    }

    openModal(video) {
        this.modalVideo.src = video.url;
        this.modalTitle.textContent = video.title;
        this.modalDescription.textContent = video.description;
        this.modal.classList.remove('hidden');
        this.modalVideo.play();
    }

    closeModal() {
        this.modal.classList.add('hidden');
        this.modalVideo.pause();
        this.modalVideo.src = '';
    }

    loadMore() {
        this.page++;
        this.loadVideos();
    }

    handleSearch() {
        this.filters.search = this.searchInput.value;
        this.resetAndReload();
    }

    handleFilter() {
        this.filters.category = this.categorySelect.value;
        this.filters.sort = this.sortSelect.value;
        this.resetAndReload();
    }

    resetAndReload() {
        this.page = 1;
        this.videoGrid.innerHTML = '';
        this.loadVideos();
    }
}

// Initialize gallery when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new VideoGallery();
}); 