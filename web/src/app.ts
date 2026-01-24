/**
 * Video recording and playback functionality
 */

class VideoRecorder {
  private stream: MediaStream | null = null;
  private mediaRecorder: MediaRecorder | null = null;
  private recordedChunks: Blob[] = [];
  private isRecording = false;

  // DOM elements
  private previewVideo: HTMLVideoElement;
  private recordBtn: HTMLButtonElement;
  private stopBtn: HTMLButtonElement;
  private playBtn: HTMLButtonElement;
  private recordingsList: HTMLElement;

  constructor() {
    const previewEl = document.getElementById('preview') as HTMLVideoElement | null;
    const recordBtnEl = document.getElementById('recordBtn') as HTMLButtonElement | null;
    const stopBtnEl = document.getElementById('stopBtn') as HTMLButtonElement | null;
    const playBtnEl = document.getElementById('playBtn') as HTMLButtonElement | null;
    const recordingsListEl = document.getElementById('recordingsList') as HTMLElement | null;

    if (!previewEl || !recordBtnEl || !stopBtnEl || !playBtnEl || !recordingsListEl) {
      throw new Error('Required DOM elements not found for VideoRecorder');
    }

    this.previewVideo = previewEl;
    this.recordBtn = recordBtnEl;
    this.stopBtn = stopBtnEl;
    this.playBtn = playBtnEl;
    this.recordingsList = recordingsListEl;

    // Bind event handlers
    this.recordBtn.addEventListener('click', () => this.startRecording());
    this.stopBtn.addEventListener('click', () => this.stopRecording());
    this.playBtn.addEventListener('click', () => this.playRecording());

    // Initialize
    this.initializeCamera();
  }

  async initializeCamera(): Promise<void> {
    try {
      this.stream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: { ideal: 1920 },
          height: { ideal: 1080 },
          facingMode: 'user',
        },
        audio: true,
      });
      this.previewVideo.srcObject = this.stream;
    } catch (err) {
      console.error('Error accessing camera:', err);
      alert('Unable to access camera. Please ensure you have granted camera permissions.');
    }
  }

  startRecording(): void {
    if (!this.stream) return;

    this.recordedChunks = [];
    this.mediaRecorder = new MediaRecorder(this.stream, {
      mimeType: 'video/webm;codecs=vp9,opus',
    });

    this.mediaRecorder.ondataavailable = (event: BlobEvent): void => {
      if (event.data.size > 0) {
        this.recordedChunks.push(event.data);
      }
    };

    this.mediaRecorder.onstop = (): void => {
      const blob = new Blob(this.recordedChunks, {
        type: 'video/webm',
      });
      this.saveRecording(blob);
    };

    this.mediaRecorder.start();
    this.isRecording = true;
    this.updateUI();
  }

  stopRecording(): void {
    if (!this.isRecording || !this.mediaRecorder) return;

    this.mediaRecorder.stop();
    this.isRecording = false;
    this.updateUI();
  }

  playRecording(): void {
    if (this.recordedChunks.length === 0) return;

    const blob = new Blob(this.recordedChunks, {
      type: 'video/webm',
    });
    const url = URL.createObjectURL(blob);
    this.previewVideo.src = url;
    this.previewVideo.play();
  }

  private saveRecording(blob: Blob): void {
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

  private updateUI(): void {
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
document.addEventListener('DOMContentLoaded', (): void => {
  try {
    new VideoRecorder();
  } catch (error) {
    console.error('Failed to initialize VideoRecorder:', error);
  }
});
