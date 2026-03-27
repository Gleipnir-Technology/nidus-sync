<style scoped>
.upload-widget {
  max-width: 400px;
  font-family: system-ui, -apple-system, sans-serif;
}

.upload-area {
  position: relative;
  border: 2px dashed #cbd5e0;
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  background-color: #f7fafc;
  transition: all 0.3s ease;
  cursor: pointer;
}

.upload-area:hover {
  border-color: #4299e1;
  background-color: #ebf8ff;
}

.file-input {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.file-input:disabled {
  cursor: not-allowed;
}

.upload-info {
  pointer-events: none;
}

.placeholder {
  color: #718096;
}

.placeholder svg {
  margin: 0 auto 1rem;
  color: #a0aec0;
}

.placeholder p {
  margin: 0;
  font-size: 0.875rem;
}

.file-selected {
  color: #2d3748;
}

.file-name {
  margin: 0 0 0.5rem;
  font-weight: 600;
  word-break: break-all;
}

.file-size {
  margin: 0;
  font-size: 0.875rem;
  color: #718096;
}

.upload-button {
  width: 100%;
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background-color: #4299e1;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.upload-button:hover {
  background-color: #3182ce;
}

.upload-button:disabled {
  background-color: #cbd5e0;
  cursor: not-allowed;
}

.upload-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #ebf8ff;
  border-radius: 6px;
  color: #2c5282;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid #bee3f8;
  border-top-color: #4299e1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.success-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #c6f6d5;
  color: #22543d;
  border-radius: 6px;
  font-weight: 500;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #fed7d7;
  color: #742a2a;
  border-radius: 6px;
  font-weight: 500;
}
</style>
<template>
  <div class="upload-widget">
    <div class="upload-area">
      <input
        ref="fileInput"
        type="file"
        accept=".csv"
        @change="handleFileSelect"
        :disabled="isUploading"
        class="file-input"
      />
      
      <div class="upload-info">
        <div v-if="!selectedFile" class="placeholder">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="17 8 12 3 7 8"></polyline>
            <line x1="12" y1="3" x2="12" y2="15"></line>
          </svg>
          <p>Select a CSV file to upload</p>
        </div>
        
        <div v-else class="file-selected">
          <p class="file-name">{{ selectedFile.name }}</p>
          <p class="file-size">{{ formatFileSize(selectedFile.size) }}</p>
        </div>
      </div>
    </div>

    <button
      v-if="selectedFile && !isUploading"
      @click="uploadFile"
      class="upload-button"
      :disabled="isUploading"
    >
      Upload File
    </button>

    <div v-if="isUploading" class="upload-status">
      <div class="spinner"></div>
      <p>Uploading...</p>
    </div>

    <div v-if="uploadSuccess" class="success-message">
      ✓ File uploaded successfully!
    </div>

    <div v-if="errorMessage" class="error-message">
      ✗ {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  uploadUrl: string;
  headers?: Record<string, string>;
  additionalData?: Record<string, string>;
}

interface Emits {
  (e: 'doError', error: Error): void;
  (e: 'doFileSelected', file: File): void;
  (e: 'doSuccess', response: any): void;
}

const props = withDefaults(defineProps<Props>(), {
  headers: () => ({}),
  additionalData: () => ({})
});

const emit = defineEmits<Emits>();

const fileInput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const isUploading = ref(false);
const uploadSuccess = ref(false);
const errorMessage = ref('');

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  
  if (file) {
    // Validate it's a CSV file
    if (!file.name.toLowerCase().endsWith('.csv')) {
      errorMessage.value = 'Please select a valid CSV file';
      selectedFile.value = null;
      return;
    }
    
    selectedFile.value = file;
    errorMessage.value = '';
    uploadSuccess.value = false;
    emit('doFileSelected', file);
  }
};

const uploadFile = async () => {
  if (!selectedFile.value) return;

  isUploading.value = true;
  errorMessage.value = '';
  uploadSuccess.value = false;

  try {
    const formData = new FormData();
    formData.append('file', selectedFile.value);
    
    // Add any additional data to the form
    Object.entries(props.additionalData).forEach(([key, value]) => {
      formData.append(key, value);
    });

    const response = await fetch(props.uploadUrl, {
      method: 'POST',
      headers: {
        ...props.headers,
        // Don't set Content-Type - let browser set it with boundary
      },
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.statusText}`);
    }

    const data = await response.json();
    uploadSuccess.value = true;
    emit('doSuccess', data);
    
    // Reset after successful upload
    setTimeout(() => {
      resetUpload();
    }, 2000);

  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Upload failed';
    emit('doError', error as Error);
  } finally {
    isUploading.value = false;
  }
};

const resetUpload = () => {
  selectedFile.value = null;
  uploadSuccess.value = false;
  errorMessage.value = '';
  if (fileInput.value) {
    fileInput.value.value = '';
  }
};

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
};

// Expose methods for parent component
defineExpose({
  resetUpload
});
</script>

