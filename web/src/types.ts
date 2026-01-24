/**
 * Shared TypeScript type definitions for Gooji frontend
 */

export interface Video {
  id: string;
  title: string;
  description: string;
  tags: string[];
  createdAt?: string;
  duration?: number;
  thumbnail?: string;
}

export interface VideoListResponse extends Array<Video> {}

export interface VideoUploadResponse {
  id: string;
  message: string;
}

export interface ApiError {
  error: string;
  message: string;
}

export type MediaRecorderState = 'inactive' | 'recording' | 'paused';

export interface MediaRecorderOptions {
  mimeType?: string;
  audioBitsPerSecond?: number;
  videoBitsPerSecond?: number;
  bitsPerSecond?: number;
}

export interface CameraConstraints {
  video: {
    width?: { ideal?: number; min?: number; max?: number };
    height?: { ideal?: number; min?: number; max?: number };
    facingMode?: 'user' | 'environment';
  };
  audio: boolean | {
    echoCancellation?: boolean;
    noiseSuppression?: boolean;
  };
}

export interface VideoElement extends HTMLVideoElement {
  srcObject: MediaStream | null;
}

export interface FormDataEntry {
  name: string;
  value: string | Blob;
  filename?: string;
}
