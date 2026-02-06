class PhotoUpload extends HTMLElement {
	// make element form-associated
	static formAssociated = true;

	constructor() {
		super();
		this.attachShadow({ mode: "open" });
		// Track all selected files
		this.selectedFiles = new Map();
		this.fileCounter = 0;
		this.render();
		this.fileInput = this.shadowRoot.getElementById("photos");
		this.internals = this.attachInternals();
	}

	connectedCallback() {
		setTimeout(() => this._initializeUploader(), 0);
	}

	_initializeUploader() {
		// Elements
		const photoInput = this.shadowRoot.querySelector("#photos");

		// Handle photo selection
		photoInput.addEventListener("change", () => {
			this._handlePhotoSelection();
			this._updateFormValue();
		});

		// Handle drag and drop
		const photoDropArea = this.shadowRoot.querySelector("#photoDropArea");

		photoDropArea.addEventListener("dragover", (e) => {
			e.preventDefault();
			photoDropArea.style.backgroundColor = "#e9ecef";
		});

		photoDropArea.addEventListener("dragleave", () => {
			photoDropArea.style.backgroundColor = "#f8f9fa";
		});

		photoDropArea.addEventListener("drop", (e) => {
			e.preventDefault();
			photoDropArea.style.backgroundColor = "#f8f9fa";

			if (e.dataTransfer.files.length) {
				this._handleFiles(e.dataTransfer.files);
			}
		});
	}

	// Update form value with all selected files
	_updateFormValue() {
		const entries = new FormData();
		for (const [fileId, file] of this.selectedFiles.entries()) {
			entries.append(`photo_${fileId}`, file);
		}
		this.internals.setFormValue(entries);
	}

	// Handle files from drag and drop
	_handleFiles(files) {
		// Set the files to the input element
		// (Not directly possible, but we can process them manually)
		Array.from(files).forEach((file) => {
			if (file.type.match("image.*")) {
				const fileId = this.fileCounter++;
				this.selectedFiles.set(fileId, file);
				this._createImagePreview(file, fileId);
			}
		});

		this._updateFormValue();
	}

	render() {
		const style = `
                <link href="/static/css/bootstrap.css" rel="stylesheet" />
                <style>
                .photo-upload-area {
                        border: 2px dashed #ccc;
                        border-radius: 8px;
                        padding: 20px;
                        text-align: center;
                        margin-bottom: 20px;
                        background-color: #f9f9f9;
                }

                .photo-preview {
                        display: flex;
                        flex-wrap: wrap;
                        gap: 10px;
                        margin-top: 15px;
                }

                .photo-preview img {
                        width: 80px;
                        height: 80px;
                        object-fit: cover;
                        border-radius: 4px;
                }
                </style>
        `;

		// Create the table
		let html = `
                <div class="photo-upload-area">
                        <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor" class="bi bi-camera mb-2" viewBox="0 0 16 16">
                                <path d="M15 12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h1.172a3 3 0 0 0 2.12-.879l.83-.828A1 1 0 0 1 6.827 3h2.344a1 1 0 0 1 .707.293l.828.828A3 3 0 0 0 12.828 5H14a1 1 0 0 1 1 1v6zM2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 9.172 2H6.828a2 2 0 0 0-1.414.586l-.828.828A2 2 0 0 1 3.172 4H2z"/>
                                <path d="M8 11a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5zm0 1a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7zM3 6.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0z"/>
                        </svg>
                        <div class="file-upload-container" id="photoDropArea">
                                <input type="file" id="photos" name="photos" class="d-none" accept="image/*" multiple>
                                <button type="button" class="btn btn-outline-primary mb-2" onclick="this.getRootNode().handleButtonClick()">Add Photos</button>
                        </div>
                        <small class="d-block text-muted">Take pictures of the mosquito problem area</small>

                        <!-- Photo Preview Area -->
                        <div id="photoPreviewContainer" class="photo-preview mt-3 d-flex flex-wrap">
                                <!-- Image previews will be added here by JavaScript -->
                        </div>
                </div>
        `;

		// Set the shadow DOM content
		this.shadowRoot.innerHTML = style + html;
		this.shadowRoot.handleButtonClick = () => {
			const photoInput = this.shadowRoot.querySelector("#photos");
			photoInput.click();
		};
	}

	/**
	 * Create an image preview for a single file
	 */
	_createImagePreview(file, fileId) {
		const photoPreviewContainer = this.shadowRoot.querySelector(
			"#photoPreviewContainer",
		);

		// Create preview container
		const previewContainer = document.createElement("div");
		previewContainer.className = "position-relative m-1";
		previewContainer.dataset.fileId = fileId;

		// Create image preview
		const img = document.createElement("img");
		img.className = "img-thumbnail";
		img.style.width = "100px";
		img.style.height = "100px";
		img.style.objectFit = "cover";

		// Read file and set preview
		const reader = new FileReader();
		reader.onload = (e) => {
			img.src = e.target.result;
		};
		reader.readAsDataURL(file);

		// Create remove button
		const removeBtn = document.createElement("button");
		removeBtn.type = "button";
		removeBtn.className = "btn btn-sm btn-danger position-absolute top-0 end-0";
		removeBtn.innerHTML = "&times;";
		removeBtn.style.fontSize = "10px";
		removeBtn.style.padding = "0 5px";

		// Handle remove button click
		removeBtn.addEventListener("click", () => {
			// Remove this file from our collection
			this.selectedFiles.delete(parseInt(previewContainer.dataset.fileId));
			// Update the form value
			this._updateFormValue();
			// Remove the preview
			previewContainer.remove();
		});

		// Add elements to the preview container
		previewContainer.appendChild(img);
		previewContainer.appendChild(removeBtn);
		photoPreviewContainer.appendChild(previewContainer);
	}

	/**
	 * Handle photo selection and preview
	 */
	_handlePhotoSelection() {
		const photoInput = this.shadowRoot.querySelector("#photos");

		// Check if files were selected
		if (photoInput.files && photoInput.files.length > 0) {
			// Loop through selected files
			Array.from(photoInput.files).forEach((file) => {
				if (!file.type.match("image.*")) {
					console.log("Skipping non-image file", file.type);
					return; // Skip non-image files
				}

				// Add file to our collection with unique ID
				const fileId = this.fileCounter++;
				this.selectedFiles.set(fileId, file);

				// Create and add preview
				this._createImagePreview(file, fileId);
			});
		}
	}
}

// Register the custom element
customElements.define("photo-upload", PhotoUpload);
