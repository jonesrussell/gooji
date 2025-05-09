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
        const stream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true
        });
        preview.srcObject = stream;
        mediaRecorder = new MediaRecorder(stream);

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

        recordBtn.disabled = false;
    } catch (err) {
        console.error('Error accessing camera:', err);
        alert('Error accessing camera. Please make sure you have granted permission.');
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
                    <img src="/api/videos/thumbnail?id=${video.id}" alt="${video.title}" class="w-full h-full object-cover">
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

// Initialize
initCamera();
loadRecordings(); 