// Video editing functionality
class VideoEditor {
    constructor() {
        this.video = document.getElementById('editor-video');
        this.trimStart = document.getElementById('trim-start');
        this.trimEnd = document.getElementById('trim-end');
        this.playbackRate = document.getElementById('playback-rate');
        this.volume = document.getElementById('volume');
        this.brightness = document.getElementById('brightness');
        this.contrast = document.getElementById('contrast');
        this.saturation = document.getElementById('saturation');

        this.initializeControls();
    }

    initializeControls() {
        // Trim controls
        this.trimStart.addEventListener('input', () => this.updateTrim());
        this.trimEnd.addEventListener('input', () => this.updateTrim());

        // Playback rate
        this.playbackRate.addEventListener('input', () => {
            this.video.playbackRate = parseFloat(this.playbackRate.value);
        });

        // Volume
        this.volume.addEventListener('input', () => {
            this.video.volume = parseFloat(this.volume.value);
        });

        // Video filters
        this.brightness.addEventListener('input', () => this.updateFilters());
        this.contrast.addEventListener('input', () => this.updateFilters());
        this.saturation.addEventListener('input', () => this.updateFilters());
    }

    updateTrim() {
        const start = parseFloat(this.trimStart.value);
        const end = parseFloat(this.trimEnd.value);

        if (this.video.currentTime < start) {
            this.video.currentTime = start;
        } else if (this.video.currentTime > end) {
            this.video.currentTime = end;
        }
    }

    updateFilters() {
        const brightness = parseFloat(this.brightness.value);
        const contrast = parseFloat(this.contrast.value);
        const saturation = parseFloat(this.saturation.value);

        this.video.style.filter = `
            brightness(${brightness})
            contrast(${contrast})
            saturate(${saturation})
        `;
    }

    // Export edited video
    async exportVideo() {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        const stream = canvas.captureStream();
        const mediaRecorder = new MediaRecorder(stream, {
            mimeType: 'video/webm;codecs=vp9,opus'
        });

        const chunks = [];
        mediaRecorder.ondataavailable = (e) => chunks.push(e.data);
        mediaRecorder.onstop = () => {
            const blob = new Blob(chunks, { type: 'video/webm' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'edited-video.webm';
            a.click();
        };

        // Start recording
        mediaRecorder.start();

        // Draw video frames
        const start = parseFloat(this.trimStart.value);
        const end = parseFloat(this.trimEnd.value);
        this.video.currentTime = start;

        const drawFrame = () => {
            if (this.video.currentTime >= end) {
                mediaRecorder.stop();
                return;
            }

            ctx.drawImage(this.video, 0, 0, canvas.width, canvas.height);
            this.video.currentTime += 1/30; // 30fps
            requestAnimationFrame(drawFrame);
        };

        drawFrame();
    }

    // Reset all controls to default values
    reset() {
        this.trimStart.value = 0;
        this.trimEnd.value = this.video.duration;
        this.playbackRate.value = 1;
        this.volume.value = 1;
        this.brightness.value = 1;
        this.contrast.value = 1;
        this.saturation.value = 1;

        this.updateTrim();
        this.updateFilters();
    }
}

// Initialize editor when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new VideoEditor();
}); 