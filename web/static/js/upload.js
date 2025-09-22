// Upload page functionality
document.addEventListener('DOMContentLoaded', function () {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('videoFile');
    const dropZoneContent = document.getElementById('dropZoneContent');
    const fileInfo = document.getElementById('fileInfo');
    const fileName = document.getElementById('fileName');
    const fileSize = document.getElementById('fileSize');
    const changeFileBtn = document.getElementById('changeFile');
    const uploadBtn = document.getElementById('uploadBtn');
    const uploadForm = document.getElementById('uploadForm');
    const uploadProgress = document.getElementById('uploadProgress');
    const progressBar = document.getElementById('progressBar');
    const uploadPercentage = document.getElementById('uploadPercentage');
    const successModal = document.getElementById('successModal');
    const successModalContent = document.getElementById('successModalContent');
    const uploadAnotherBtn = document.getElementById('uploadAnother');

    let selectedFile = null;

    // File size formatting
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    // Validate file type
    function isValidVideoFile(file) {
        const allowedTypes = ['video/mp4', 'video/webm', 'video/quicktime', 'video/x-msvideo'];
        return allowedTypes.includes(file.type);
    }

    // Validate file size (500MB limit)
    function isValidFileSize(file) {
        const maxSize = 500 * 1024 * 1024; // 500MB in bytes
        return file.size <= maxSize;
    }

    // Show file info
    function showFileInfo(file) {
        fileName.textContent = file.name;
        fileSize.textContent = formatFileSize(file.size);
        dropZoneContent.classList.add('hidden');
        fileInfo.classList.remove('hidden');
        uploadBtn.disabled = false;
        dropZone.classList.remove('border-gray-300', 'hover:border-indigo-400', 'hover:bg-indigo-50');
        dropZone.classList.add('border-green-400', 'bg-green-50');
    }

    // Reset file selection
    function resetFileSelection() {
        selectedFile = null;
        dropZoneContent.classList.remove('hidden');
        fileInfo.classList.add('hidden');
        uploadBtn.disabled = true;
        dropZone.classList.remove('border-green-400', 'bg-green-50');
        dropZone.classList.add('border-gray-300', 'hover:border-indigo-400', 'hover:bg-indigo-50');
        fileInput.value = '';
    }

    // Handle file selection
    function handleFileSelect(file) {
        if (!isValidVideoFile(file)) {
            alert('Please select a valid video file (MP4, WebM, MOV, AVI).');
            return;
        }

        if (!isValidFileSize(file)) {
            alert('File size must be less than 500MB.');
            return;
        }

        selectedFile = file;
        showFileInfo(file);
    }

    // Drag and drop handlers
    dropZone.addEventListener('dragover', function (e) {
        e.preventDefault();
        dropZone.classList.add('border-indigo-400', 'bg-indigo-50');
    });

    dropZone.addEventListener('dragleave', function (e) {
        e.preventDefault();
        dropZone.classList.remove('border-indigo-400', 'bg-indigo-50');
    });

    dropZone.addEventListener('drop', function (e) {
        e.preventDefault();
        dropZone.classList.remove('border-indigo-400', 'bg-indigo-50');

        const files = e.dataTransfer.files;
        if (files.length > 0) {
            handleFileSelect(files[0]);
        }
    });

    // File input click handler - prevent reopening if file already selected
    fileInput.addEventListener('click', function (e) {
        if (selectedFile) {
            e.preventDefault();
            e.stopPropagation();
            return false;
        }
    });

    // File input change
    fileInput.addEventListener('change', function (e) {
        const file = e.target.files[0];
        if (file) {
            handleFileSelect(file);
        }
    });

    // Change file button
    changeFileBtn.addEventListener('click', function (e) {
        e.stopPropagation();
        resetFileSelection();
    });

    // Upload form submission
    uploadForm.addEventListener('submit', async function (e) {
        e.preventDefault();

        if (!selectedFile) {
            alert('Please select a video file first.');
            return;
        }

        const formData = new FormData();
        formData.append('video', selectedFile);
        formData.append('title', document.getElementById('title').value);
        formData.append('description', document.getElementById('description').value);
        formData.append('category', document.getElementById('category').value);
        formData.append('tags', document.getElementById('tags').value);
        formData.append('language', document.getElementById('language').value);
        formData.append('public', document.getElementById('public').checked);

        // Show progress bar
        uploadProgress.classList.remove('hidden');
        uploadBtn.disabled = true;
        uploadBtn.innerHTML = `
            <div class="flex items-center justify-center space-x-3">
                <div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                <span>Uploading...</span>
            </div>
        `;

        try {
            const xhr = new XMLHttpRequest();

            // Upload progress
            xhr.upload.addEventListener('progress', function (e) {
                if (e.lengthComputable) {
                    const percentComplete = (e.loaded / e.total) * 100;
                    progressBar.style.width = percentComplete + '%';
                    uploadPercentage.textContent = Math.round(percentComplete) + '%';
                }
            });

            // Upload complete
            xhr.addEventListener('load', function () {
                if (xhr.status === 200) {
                    // Show success modal
                    showSuccessModal();
                } else {
                    alert('Upload failed. Please try again.');
                    resetUploadState();
                }
            });

            // Upload error
            xhr.addEventListener('error', function () {
                alert('Upload failed. Please check your connection and try again.');
                resetUploadState();
            });

            xhr.open('POST', '/api/videos');
            xhr.send(formData);

        } catch (error) {
            console.error('Upload error:', error);
            alert('Upload failed. Please try again.');
            resetUploadState();
        }
    });

    // Reset upload state
    function resetUploadState() {
        uploadProgress.classList.add('hidden');
        uploadBtn.disabled = false;
        uploadBtn.innerHTML = `
            <div class="flex items-center justify-center space-x-3">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12">
                    </path>
                </svg>
                <span>Upload Video</span>
            </div>
        `;
        progressBar.style.width = '0%';
        uploadPercentage.textContent = '0%';
    }

    // Show success modal
    function showSuccessModal() {
        successModal.classList.remove('hidden');
        setTimeout(() => {
            successModalContent.classList.remove('scale-95', 'opacity-0');
            successModalContent.classList.add('scale-100', 'opacity-100');
        }, 10);
    }

    // Hide success modal
    function hideSuccessModal() {
        successModalContent.classList.remove('scale-100', 'opacity-100');
        successModalContent.classList.add('scale-95', 'opacity-0');
        setTimeout(() => {
            successModal.classList.add('hidden');
        }, 200);
    }

    // Upload another button
    uploadAnotherBtn.addEventListener('click', function () {
        hideSuccessModal();
        resetFileSelection();
        uploadForm.reset();
        resetUploadState();

        // Scroll to top
        window.scrollTo({ top: 0, behavior: 'smooth' });
    });

    // Close modal when clicking outside
    successModal.addEventListener('click', function (e) {
        if (e.target === successModal) {
            hideSuccessModal();
        }
    });

    // Close modal with escape key
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape' && !successModal.classList.contains('hidden')) {
            hideSuccessModal();
        }
    });

    // Form validation
    const requiredFields = ['title', 'description', 'category'];
    requiredFields.forEach(fieldId => {
        const field = document.getElementById(fieldId);
        field.addEventListener('input', function () {
            const isValid = requiredFields.every(id => {
                const element = document.getElementById(id);
                return element.value.trim() !== '';
            });

            if (selectedFile && isValid) {
                uploadBtn.disabled = false;
            } else {
                uploadBtn.disabled = true;
            }
        });
    });


    // Auto-resize textarea
    const descriptionTextarea = document.getElementById('description');
    descriptionTextarea.addEventListener('input', function () {
        this.style.height = 'auto';
        this.style.height = this.scrollHeight + 'px';
    });

    // Tag input enhancement
    const tagsInput = document.getElementById('tags');
    tagsInput.addEventListener('input', function () {
        // Remove extra spaces and normalize
        this.value = this.value.replace(/\s+/g, ' ').trim();
    });

    // Category change handler
    const categorySelect = document.getElementById('category');
    categorySelect.addEventListener('change', function () {
        // Auto-suggest tags based on category
        const tagSuggestions = {
            'language': 'ojibwe, language, learning, pronunciation',
            'culture': 'tradition, culture, ceremony, customs',
            'story': 'story, legend, tale, narrative',
            'history': 'history, historical, past, ancestors',
            'music': 'music, song, singing, melody',
            'crafts': 'craft, traditional, making, art',
            'ceremony': 'ceremony, ritual, sacred, spiritual'
        };

        const suggestion = tagSuggestions[this.value];
        if (suggestion && !tagsInput.value.trim()) {
            tagsInput.placeholder = `Suggested: ${suggestion}`;
        }
    });
});

