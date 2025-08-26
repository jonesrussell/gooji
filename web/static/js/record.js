// Video recording elements
const preview = document.getElementById('preview');
const recordBtn = document.getElementById('recordBtn');
const stopBtn = document.getElementById('stopBtn');
const uploadBtn = document.getElementById('uploadBtn');
const uploadForm = document.getElementById('uploadForm');
const recordingsList = document.getElementById('recordingsList');

let mediaRecorder;
let recordedChunks = [];

// Initialize camera
async function initCamera() {
    try {
        console.log('Initializing camera...');

        // Check if getUserMedia is supported
        if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
            throw new Error('getUserMedia is not supported in this browser');
        }

        // Check available devices
        const devices = await navigator.mediaDevices.enumerateDevices();
        const videoDevices = devices.filter(device => device.kind === 'videoinput');
        console.log('Available video devices:', videoDevices);

        if (videoDevices.length === 0) {
            throw new Error('No video devices found');
        }

        const stream = await navigator.mediaDevices.getUserMedia({
            video: {
                width: { ideal: 1280 },
                height: { ideal: 720 },
                facingMode: 'user'
            },
            audio: true
        });

        console.log('Camera stream obtained successfully');
        preview.srcObject = stream;

        // Wait for video to be ready
        preview.onloadedmetadata = () => {
            console.log('Video metadata loaded');
            recordBtn.disabled = false;
        };

        // Try to use preferred MIME type, fallback to default
        let mimeType = 'video/webm;codecs=vp9,opus';
        if (!MediaRecorder.isTypeSupported(mimeType)) {
            mimeType = 'video/webm';
            if (!MediaRecorder.isTypeSupported(mimeType)) {
                mimeType = '';
            }
        }

        mediaRecorder = new MediaRecorder(stream, mimeType ? { mimeType } : {});

        mediaRecorder.ondataavailable = (event) => {
            if (event.data.size > 0) {
                recordedChunks.push(event.data);
            }
        };

        mediaRecorder.onstop = () => {
            const blob = new Blob(recordedChunks, { type: 'video/webm' });
            const url = URL.createObjectURL(blob);
            preview.srcObject = null;
            preview.src = url;
            uploadBtn.disabled = false;
        };

    } catch (err) {
        console.error('Error accessing camera:', err);

        // Provide more specific error messages
        let errorMessage = 'Error accessing camera. ';
        if (err.name === 'NotAllowedError') {
            errorMessage += 'Please grant camera and microphone permissions and refresh the page.';
        } else if (err.name === 'NotFoundError') {
            errorMessage += 'No camera found. Please check your device connections.';
        } else if (err.name === 'NotSupportedError') {
            errorMessage += 'Your browser does not support video recording.';
        } else {
            errorMessage += err.message;
        }

        alert(errorMessage);

        // Show error in UI
        const errorDiv = document.createElement('div');
        errorDiv.className = 'bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4';
        errorDiv.innerHTML = `
            <strong>Camera Error:</strong> ${errorMessage}
            <br><small>Please check your browser permissions and refresh the page.</small>
        `;
        preview.parentElement.insertBefore(errorDiv, preview);
    }
}

// Start recording
recordBtn.addEventListener('click', () => {
    recordedChunks = [];
    mediaRecorder.start();
    recordBtn.disabled = true;
    stopBtn.disabled = false;
    uploadBtn.disabled = true;
});

// Stop recording
stopBtn.addEventListener('click', () => {
    mediaRecorder.stop();
    recordBtn.disabled = false;
    stopBtn.disabled = true;
});

// Handle form submission
uploadForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const formData = new FormData();
    const blob = new Blob(recordedChunks, { type: 'video/webm' });
    formData.append('video', blob, 'recording.webm');
    formData.append('title', document.getElementById('title').value);
    formData.append('description', document.getElementById('description').value);
    formData.append('tags', document.getElementById('tags').value);

    try {
        const response = await fetch('/api/videos', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Upload failed');
        }

        const result = await response.json();
        alert('Video uploaded successfully!');

        // Reset form and recording
        uploadForm.reset();
        recordedChunks = [];
        preview.srcObject = null;
        uploadBtn.disabled = true;

        // Refresh recordings list
        loadRecordings();
    } catch (err) {
        console.error('Error uploading video:', err);
        alert('Error uploading video. Please try again.');
    }
});

// Load recent recordings
async function loadRecordings() {
    try {
        const response = await fetch('/api/videos');
        if (!response.ok) {
            throw new Error('Failed to load recordings');
        }

        const videos = await response.json();
        recordingsList.innerHTML = '';

        videos.forEach(video => {
            const videoCard = document.createElement('div');
            videoCard.className = 'bg-gray-100 rounded-lg overflow-hidden';
            videoCard.innerHTML = `
                <div class="aspect-w-16 aspect-h-9">
                    <img src="/api/thumbnails?id=${video.id}" alt="${video.title}" class="w-full h-full object-cover">
                </div>
                <div class="p-4">
                    <h3 class="font-semibold text-gray-800">${video.title}</h3>
                    <p class="text-sm text-gray-600">${video.description}</p>
                    <div class="mt-2 flex flex-wrap gap-2">
                        ${video.tags.map(tag => `
                            <span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">${tag}</span>
                        `).join('')}
                    </div>
                </div>
            `;
            recordingsList.appendChild(videoCard);
        });
    } catch (err) {
        console.error('Error loading recordings:', err);
    }
}

// Initialize with retry mechanism
let cameraRetryCount = 0;
const maxRetries = 3;

async function initializeWithRetry() {
    try {
        await initCamera();
        loadRecordings();
    } catch (err) {
        console.error('Camera initialization failed:', err);
        if (cameraRetryCount < maxRetries) {
            cameraRetryCount++;
            console.log(`Retrying camera initialization (${cameraRetryCount}/${maxRetries})...`);
            setTimeout(initializeWithRetry, 1000);
        } else {
            console.error('Max retries reached for camera initialization');
        }
    }
}

initializeWithRetry();
