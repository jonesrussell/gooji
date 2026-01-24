/**
 * Video editing functionality
 */

class VideoEditor {
  private video: HTMLVideoElement;
  private trimStart: HTMLInputElement;
  private trimEnd: HTMLInputElement;
  private playbackRate: HTMLInputElement;
  private volume: HTMLInputElement;
  private brightness: HTMLInputElement;
  private contrast: HTMLInputElement;
  private saturation: HTMLInputElement;

  constructor() {
    const videoEl = document.getElementById('editor-video') as HTMLVideoElement | null;
    const trimStartEl = document.getElementById('trim-start') as HTMLInputElement | null;
    const trimEndEl = document.getElementById('trim-end') as HTMLInputElement | null;
    const playbackRateEl = document.getElementById('playback-rate') as HTMLInputElement | null;
    const volumeEl = document.getElementById('volume') as HTMLInputElement | null;
    const brightnessEl = document.getElementById('brightness') as HTMLInputElement | null;
    const contrastEl = document.getElementById('contrast') as HTMLInputElement | null;
    const saturationEl = document.getElementById('saturation') as HTMLInputElement | null;

    if (!videoEl || !trimStartEl || !trimEndEl || !playbackRateEl || !volumeEl || 
        !brightnessEl || !contrastEl || !saturationEl) {
      throw new Error('Required DOM elements not found for VideoEditor');
    }

    this.video = videoEl;
    this.trimStart = trimStartEl;
    this.trimEnd = trimEndEl;
    this.playbackRate = playbackRateEl;
    this.volume = volumeEl;
    this.brightness = brightnessEl;
    this.contrast = contrastEl;
    this.saturation = saturationEl;

    this.initializeControls();
  }

  private initializeControls(): void {
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

  private updateTrim(): void {
    const start = parseFloat(this.trimStart.value);
    const end = parseFloat(this.trimEnd.value);

    if (this.video.currentTime < start) {
      this.video.currentTime = start;
    } else if (this.video.currentTime > end) {
      this.video.currentTime = end;
    }
  }

  private updateFilters(): void {
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
  async exportVideo(): Promise<void> {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      throw new Error('Could not get canvas context');
    }

    canvas.width = this.video.videoWidth;
    canvas.height = this.video.videoHeight;

    const stream = canvas.captureStream();
    const mediaRecorder = new MediaRecorder(stream, {
      mimeType: 'video/webm;codecs=vp9,opus',
    });

    const chunks: Blob[] = [];
    mediaRecorder.ondataavailable = (e: BlobEvent): void => {
      if (e.data.size > 0) {
        chunks.push(e.data);
      }
    };

    mediaRecorder.onstop = (): void => {
      const blob = new Blob(chunks, { type: 'video/webm' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'edited-video.webm';
      a.click();
      URL.revokeObjectURL(url);
    };

    // Start recording
    mediaRecorder.start();

    // Draw video frames
    const start = parseFloat(this.trimStart.value);
    const end = parseFloat(this.trimEnd.value);
    this.video.currentTime = start;

    const drawFrame = (): void => {
      if (this.video.currentTime >= end) {
        mediaRecorder.stop();
        return;
      }

      ctx.drawImage(this.video, 0, 0, canvas.width, canvas.height);
      this.video.currentTime += 1 / 30; // 30fps
      requestAnimationFrame(drawFrame);
    };

    drawFrame();
  }

  // Reset all controls to default values
  reset(): void {
    this.trimStart.value = '0';
    this.trimEnd.value = this.video.duration.toString();
    this.playbackRate.value = '1';
    this.volume.value = '1';
    this.brightness.value = '1';
    this.contrast.value = '1';
    this.saturation.value = '1';

    this.updateTrim();
    this.updateFilters();
  }
}

// Initialize editor when the page loads
document.addEventListener('DOMContentLoaded', (): void => {
  try {
    new VideoEditor();
  } catch (error) {
    console.error('Failed to initialize VideoEditor:', error);
  }
});
