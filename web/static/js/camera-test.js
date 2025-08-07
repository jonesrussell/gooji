// Camera test functionality
class CameraTest {
    constructor() {
        console.log('CameraTest constructor called');
        
        this.video = document.getElementById('testVideo');
        this.status = document.getElementById('cameraStatus');
        this.deviceList = document.getElementById('deviceList');
        this.retryBtn = document.getElementById('retryBtn');
        this.permissionsBtn = document.getElementById('permissionsBtn');

        console.log('Camera test elements found:', {
            video: !!this.video,
            status: !!this.status,
            deviceList: !!this.deviceList,
            retryBtn: !!this.retryBtn,
            permissionsBtn: !!this.permissionsBtn
        });

        // Only initialize if we're on the camera test page
        if (!this.video || !this.status) {
            console.log('Camera test elements not found, skipping initialization');
            return;
        }

        console.log('Initializing camera test functionality');

        this.retryBtn.addEventListener('click', () => this.initializeCamera());
        this.permissionsBtn.addEventListener('click', () => this.checkPermissions());

        this.initializeCamera();
        this.loadDeviceInfo();
    }

    updateStatus(message, type = 'info') {
        const colors = {
            info: 'bg-blue-100 border-blue-400 text-blue-700',
            success: 'bg-green-100 border-green-400 text-green-700',
            error: 'bg-red-100 border-red-400 text-red-700',
            warning: 'bg-yellow-100 border-yellow-400 text-yellow-700'
        };

        this.status.innerHTML = `
            <div class="border px-4 py-3 rounded ${colors[type]}">
                <strong>Status:</strong> ${message}
            </div>
        `;
    }

    async loadDeviceInfo() {
        try {
            const devices = await navigator.mediaDevices.enumerateDevices();
            const videoDevices = devices.filter(device => device.kind === 'videoinput');
            const audioDevices = devices.filter(device => device.kind === 'audioinput');

            let html = '<div class="space-y-2">';
            
            html += '<h4 class="font-semibold">Video Devices:</h4>';
            if (videoDevices.length === 0) {
                html += '<p class="text-red-600">No video devices found</p>';
            } else {
                videoDevices.forEach((device, index) => {
                    html += `<div class="ml-4">${index + 1}. ${device.label || 'Unknown device'}</div>`;
                });
            }

            html += '<h4 class="font-semibold mt-4">Audio Devices:</h4>';
            if (audioDevices.length === 0) {
                html += '<p class="text-red-600">No audio devices found</p>';
            } else {
                audioDevices.forEach((device, index) => {
                    html += `<div class="ml-4">${index + 1}. ${device.label || 'Unknown device'}</div>`;
                });
            }

            html += '</div>';
            this.deviceList.innerHTML = html;
        } catch (err) {
            this.deviceList.innerHTML = `<div class="text-red-600">Error loading device info: ${err.message}</div>`;
        }
    }

    async checkPermissions() {
        try {
            const permissions = await navigator.permissions.query({ name: 'camera' });
            const audioPermissions = await navigator.permissions.query({ name: 'microphone' });
            
            let message = `Camera: ${permissions.state}, Microphone: ${audioPermissions.state}`;
            let type = 'info';
            
            if (permissions.state === 'denied' || audioPermissions.state === 'denied') {
                type = 'error';
                message += ' - Please grant permissions in your browser settings';
            } else if (permissions.state === 'granted' && audioPermissions.state === 'granted') {
                type = 'success';
                message += ' - Permissions granted';
            }
            
            this.updateStatus(message, type);
        } catch (err) {
            this.updateStatus(`Permission check failed: ${err.message}`, 'error');
        }
    }

    async initializeCamera() {
        this.updateStatus('Initializing camera...', 'info');

        try {
            // Check if getUserMedia is supported
            if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
                throw new Error('getUserMedia is not supported in this browser');
            }

            const stream = await navigator.mediaDevices.getUserMedia({
                video: {
                    width: { ideal: 1280 },
                    height: { ideal: 720 },
                    facingMode: 'user'
                },
                audio: true
            });

            this.video.srcObject = stream;
            
            this.video.onloadedmetadata = () => {
                this.updateStatus('Camera initialized successfully!', 'success');
            };

            this.video.onerror = (err) => {
                this.updateStatus(`Video error: ${err.message}`, 'error');
            };

        } catch (err) {
            console.error('Camera initialization error:', err);
            
            let errorMessage = 'Camera initialization failed. ';
            if (err.name === 'NotAllowedError') {
                errorMessage += 'Please grant camera and microphone permissions.';
            } else if (err.name === 'NotFoundError') {
                errorMessage += 'No camera found. Please check your device connections.';
            } else if (err.name === 'NotSupportedError') {
                errorMessage += 'Your browser does not support video recording.';
            } else {
                errorMessage += err.message;
            }
            
            this.updateStatus(errorMessage, 'error');
        }
    }
}

// Initialize when page loads
document.addEventListener('DOMContentLoaded', () => {
    new CameraTest();
});
