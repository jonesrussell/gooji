/**
 * Upload page functionality
 */

document.addEventListener('DOMContentLoaded', (): void => {
  const dropZone = document.getElementById('dropZone') as HTMLElement | null;
  const fileInput = document.getElementById('videoFile') as HTMLInputElement | null;
  const dropZoneContent = document.getElementById('dropZoneContent') as HTMLElement | null;
  const fileInfo = document.getElementById('fileInfo') as HTMLElement | null;
  const fileName = document.getElementById('fileName') as HTMLElement | null;
  const fileSize = document.getElementById('fileSize') as HTMLElement | null;
  const changeFileBtn = document.getElementById('changeFile') as HTMLButtonElement | null;
  const uploadBtn = document.getElementById('uploadBtn') as HTMLButtonElement | null;
  const uploadForm = document.getElementById('uploadForm') as HTMLFormElement | null;
  const uploadProgress = document.getElementById('uploadProgress') as HTMLElement | null;
  const progressBar = document.getElementById('progressBar') as HTMLElement | null;
  const uploadPercentage = document.getElementById('uploadPercentage') as HTMLElement | null;
  const successModal = document.getElementById('successModal') as HTMLElement | null;
  const successModalContent = document.getElementById('successModalContent') as HTMLElement | null;
  const uploadAnotherBtn = document.getElementById('uploadAnother') as HTMLButtonElement | null;

  let selectedFile: File | null = null;

  if (!dropZone || !fileInput || !dropZoneContent || !fileInfo || !fileName || !fileSize || 
      !changeFileBtn || !uploadBtn || !uploadForm || !uploadProgress || !progressBar || 
      !uploadPercentage || !successModal || !successModalContent || !uploadAnotherBtn) {
    console.error('Required DOM elements not found');
    return;
  }

  // File size formatting
  function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  }

  // Validate file type
  function isValidVideoFile(file: File): boolean {
    const allowedTypes = ['video/mp4', 'video/webm', 'video/quicktime', 'video/x-msvideo'];
    return allowedTypes.includes(file.type);
  }

  // Validate file size (500MB limit)
  function isValidFileSize(file: File): boolean {
    const maxSize = 500 * 1024 * 1024; // 500MB in bytes
    return file.size <= maxSize;
  }

  // Show file info
  function showFileInfo(file: File): void {
    fileName.textContent = file.name;
    fileSize.textContent = formatFileSize(file.size);
    dropZoneContent.classList.add('hidden');
    fileInfo.classList.remove('hidden');
    uploadBtn.disabled = false;
    dropZone.classList.remove('border-gray-300', 'hover:border-yellow-400', 'hover:bg-yellow-50');
    dropZone.classList.add('border-green-400', 'bg-green-50');
  }

  // Reset file selection
  function resetFileSelection(): void {
    selectedFile = null;
    dropZoneContent.classList.remove('hidden');
    fileInfo.classList.add('hidden');
    uploadBtn.disabled = true;
    dropZone.classList.remove('border-green-400', 'bg-green-50');
    dropZone.classList.add('border-gray-300', 'hover:border-yellow-400', 'hover:bg-yellow-50');
    fileInput.value = '';
  }

  // Handle file selection
  function handleFileSelect(file: File): void {
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
  dropZone.addEventListener('dragover', (e: DragEvent): void => {
    e.preventDefault();
    dropZone.classList.add('border-yellow-400', 'bg-yellow-50');
  });

  dropZone.addEventListener('dragleave', (e: DragEvent): void => {
    e.preventDefault();
    dropZone.classList.remove('border-yellow-400', 'bg-yellow-50');
  });

  dropZone.addEventListener('drop', (e: DragEvent): void => {
    e.preventDefault();
    dropZone.classList.remove('border-yellow-400', 'bg-yellow-50');

    const files = e.dataTransfer?.files;
    if (files && files.length > 0) {
      handleFileSelect(files[0]);
    }
  });

  // File input click handler - prevent reopening if file already selected
  fileInput.addEventListener('click', (e: MouseEvent): void => {
    if (selectedFile) {
      e.preventDefault();
      e.stopPropagation();
    }
  });

  // File input change
  fileInput.addEventListener('change', (e: Event): void => {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) {
      handleFileSelect(file);
    }
  });

  // Change file button
  changeFileBtn.addEventListener('click', (e: MouseEvent): void => {
    e.stopPropagation();
    resetFileSelection();
  });

  // Reset upload state
  function resetUploadState(): void {
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
    if (progressBar) {
      progressBar.style.width = '0%';
    }
    if (uploadPercentage) {
      uploadPercentage.textContent = '0%';
    }
  }

  // Upload form submission
  uploadForm.addEventListener('submit', async (e: Event): Promise<void> => {
    e.preventDefault();

    if (!selectedFile) {
      alert('Please select a video file first.');
      return;
    }

    const titleEl = document.getElementById('title') as HTMLInputElement | null;
    const descriptionEl = document.getElementById('description') as HTMLTextAreaElement | null;
    const categoryEl = document.getElementById('category') as HTMLSelectElement | null;
    const tagsEl = document.getElementById('tags') as HTMLInputElement | null;
    const languageEl = document.getElementById('language') as HTMLSelectElement | null;
    const publicEl = document.getElementById('public') as HTMLInputElement | null;

    if (!titleEl || !descriptionEl || !categoryEl || !tagsEl || !languageEl || !publicEl) {
      alert('Form fields not found');
      return;
    }

    const formData = new FormData();
    formData.append('video', selectedFile);
    formData.append('title', titleEl.value);
    formData.append('description', descriptionEl.value);
    formData.append('category', categoryEl.value);
    formData.append('tags', tagsEl.value);
    formData.append('language', languageEl.value);
    formData.append('public', publicEl.checked.toString());

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
      xhr.upload.addEventListener('progress', (e: ProgressEvent): void => {
        if (e.lengthComputable) {
          const percentComplete = (e.loaded / e.total) * 100;
          if (progressBar) {
            progressBar.style.width = `${percentComplete}%`;
          }
          if (uploadPercentage) {
            uploadPercentage.textContent = `${Math.round(percentComplete)}%`;
          }
        }
      });

      // Upload complete
      xhr.addEventListener('load', (): void => {
        if (xhr.status === 200) {
          // Show success modal
          showSuccessModal();
        } else {
          alert('Upload failed. Please try again.');
          resetUploadState();
        }
      });

      // Upload error
      xhr.addEventListener('error', (): void => {
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

  // Show success modal
  function showSuccessModal(): void {
    successModal.classList.remove('hidden');
    setTimeout(() => {
      successModalContent.classList.remove('scale-95', 'opacity-0');
      successModalContent.classList.add('scale-100', 'opacity-100');
    }, 10);
  }

  // Hide success modal
  function hideSuccessModal(): void {
    successModalContent.classList.remove('scale-100', 'opacity-100');
    successModalContent.classList.add('scale-95', 'opacity-0');
    setTimeout(() => {
      successModal.classList.add('hidden');
    }, 200);
  }

  // Upload another button
  uploadAnotherBtn.addEventListener('click', (): void => {
    hideSuccessModal();
    resetFileSelection();
    uploadForm.reset();
    resetUploadState();

    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });
  });

  // Close modal when clicking outside
  successModal.addEventListener('click', (e: MouseEvent): void => {
    if (e.target === successModal) {
      hideSuccessModal();
    }
  });

  // Close modal with escape key
  document.addEventListener('keydown', (e: KeyboardEvent): void => {
    if (e.key === 'Escape' && !successModal.classList.contains('hidden')) {
      hideSuccessModal();
    }
  });

  // Form validation
  const requiredFields = ['title', 'description', 'category'];
  requiredFields.forEach(fieldId => {
    const field = document.getElementById(fieldId) as HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement | null;
    if (field) {
      field.addEventListener('input', (): void => {
        const isValid = requiredFields.every(id => {
          const element = document.getElementById(id) as HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement | null;
          return element && element.value.trim() !== '';
        });

        if (selectedFile && isValid) {
          uploadBtn.disabled = false;
        } else {
          uploadBtn.disabled = true;
        }
      });
    }
  });

  // Auto-resize textarea
  const descriptionTextarea = document.getElementById('description') as HTMLTextAreaElement | null;
  if (descriptionTextarea) {
    descriptionTextarea.addEventListener('input', function (this: HTMLTextAreaElement): void {
      this.style.height = 'auto';
      this.style.height = `${this.scrollHeight}px`;
    });
  }

  // Tag input enhancement
  const tagsInput = document.getElementById('tags') as HTMLInputElement | null;
  if (tagsInput) {
    tagsInput.addEventListener('input', function (this: HTMLInputElement): void {
      // Remove extra spaces and normalize
      this.value = this.value.replace(/\s+/g, ' ').trim();
    });
  }

  // Category change handler
  const categorySelect = document.getElementById('category') as HTMLSelectElement | null;
  if (categorySelect && tagsInput) {
    categorySelect.addEventListener('change', function (this: HTMLSelectElement): void {
      // Auto-suggest tags based on category
      const tagSuggestions: Record<string, string> = {
        'language': 'ojibwe, language, learning, pronunciation',
        'culture': 'tradition, culture, ceremony, customs',
        'story': 'story, legend, tale, narrative',
        'history': 'history, historical, past, ancestors',
        'music': 'music, song, singing, melody',
        'crafts': 'craft, traditional, making, art',
        'ceremony': 'ceremony, ritual, sacred, spiritual',
      };

      const suggestion = tagSuggestions[this.value];
      if (suggestion && !tagsInput.value.trim()) {
        tagsInput.placeholder = `Suggested: ${suggestion}`;
      }
    });
  }
});
