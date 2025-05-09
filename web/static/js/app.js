// Video recording and playback functionality
class VideoRecorder {
    constructor() {
        this.stream = null;
        this.mediaRecorder = null;
        this.recordedChunks = [];
        this.isRecording = false;

        // DOM elements
        this.previewVideo = document.getElementById('preview');
        this.recordBtn = document.getElementById('recordBtn');
        this.stopBtn = document.getElementById('stopBtn');
        this.playBtn = document.getElementById('playBtn');
        this.recordingsList = document.getElementById('recordingsList');

        // Bind event handlers
        this.recordBtn.addEventListener('click', () => this.startRecording());
        this.stopBtn.addEventListener('click', () => this.stopRecording());
        this.playBtn.addEventListener('click', () => this.playRecording());

        // Initialize
        this.initializeCamera();
    }

    async initializeCamera() {
        try {
            this.stream = await navigator.mediaDevices.getUserMedia({
                video: {
                    width: { ideal: 1920 },
                    height: { ideal: 1080 },
                    facingMode: 'user'
                },
                audio: true
            });
            this.previewVideo.srcObject = this.stream;
        } catch (err) {
            console.error('Error accessing camera:', err);
            alert('Unable to access camera. Please ensure you have granted camera permissions.');
        }
    }

    startRecording() {
        if (!this.stream) return;

        this.recordedChunks = [];
        this.mediaRecorder = new MediaRecorder(this.stream, {
            mimeType: 'video/webm;codecs=vp9,opus'
        });

        this.mediaRecorder.ondataavailable = (event) => {
            if (event.data.size > 0) {
                this.recordedChunks.push(event.data);
            }
        };

        this.mediaRecorder.onstop = () => {
            const blob = new Blob(this.recordedChunks, {
                type: 'video/webm'
            });
            this.saveRecording(blob);
        };

        this.mediaRecorder.start();
        this.isRecording = true;
        this.updateUI();
    }

    stopRecording() {
        if (!this.isRecording) return;

        this.mediaRecorder.stop();
        this.isRecording = false;
        this.updateUI();
    }

    playRecording() {
        if (this.recordedChunks.length === 0) return;

        const blob = new Blob(this.recordedChunks, {
            type: 'video/webm'
        });
        const url = URL.createObjectURL(blob);
        this.previewVideo.src = url;
        this.previewVideo.play();
    }

    saveRecording(blob) {
        const url = URL.createObjectURL(blob);
        const timestamp = new Date().toISOString();
        
        // Create thumbnail
        const thumbnail = document.createElement('div');
        thumbnail.className = 'recording-thumbnail';
        
        const video = document.createElement('video');
        video.src = url;
        video.muted = true;
        
        const controls = document.createElement('div');
        controls.className = 'absolute bottom-0 left-0 right-0 bg-black bg-opacity-50 text-white p-2';
        controls.textContent = new Date(timestamp).toLocaleString();
        
        thumbnail.appendChild(video);
        thumbnail.appendChild(controls);
        
        // Add click handler to play the recording
        thumbnail.addEventListener('click', () => {
            this.previewVideo.src = url;
            this.previewVideo.play();
        });
        
        this.recordingsList.insertBefore(thumbnail, this.recordingsList.firstChild);
    }

    updateUI() {
        this.recordBtn.disabled = this.isRecording;
        this.stopBtn.disabled = !this.isRecording;
        this.playBtn.disabled = this.isRecording || this.recordedChunks.length === 0;
        
        if (this.isRecording) {
            this.recordBtn.classList.add('recording');
        } else {
            this.recordBtn.classList.remove('recording');
        }
    }
}

// Initialize the video recorder when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new VideoRecorder();
}); 