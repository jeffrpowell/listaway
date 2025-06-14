require("../navbar");

// Store collection data
let collectionId = null;
let collectionName = null;

// Event handler when the page loads
document.addEventListener('DOMContentLoaded', () => {
  // Get collection ID and name from the form
  const idInput = document.getElementById('collectionId');
  const nameInput = document.getElementById('name');
  
  if (idInput) collectionId = idInput.value;
  if (nameInput) collectionName = nameInput.value;
  
  // Initialize form handling
  initializeForm();
  initializeDeleteModal();
});

// Initialize form behavior
function initializeForm() {
  const form = document.getElementById('editCollectionForm');
  if (form) {
    form.addEventListener('submit', handleFormSubmit);
  }
}

// Initialize delete modal
function initializeDeleteModal() {
  const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
  if (confirmDeleteBtn) {
    confirmDeleteBtn.addEventListener('click', () => {
      const input = document.getElementById('confirmCollectionName');
      if (input && input.value === collectionName) {
        deleteCollection(collectionName);
      } else {
        // Show error message for incorrect name
        input.classList.add('input-error');
        setTimeout(() => {
          input.classList.remove('input-error');
        }, 1500);
      }
    });
  }
}

// Handle form submission
function handleFormSubmit(event) {
  event.preventDefault();
  
  // Get form data
  const nameInput = document.getElementById('name');
  const descriptionInput = document.getElementById('description');
  
  // Validate inputs
  if (!nameInput || !nameInput.value.trim()) {
    showInputError(nameInput, 'Collection name is required');
    return;
  }
  
  // Prepare collection data
  const collectionData = {
    name: nameInput.value.trim(),
    description: descriptionInput ? descriptionInput.value.trim() : ''
  };
  
  // Submit the data
  updateCollection(collectionData);
}

// Show error on input
function showInputError(inputElement, message) {
  if (inputElement) {
    inputElement.classList.add('input-error');
    
    // Add error message if not already present
    let errorMessage = inputElement.nextElementSibling;
    if (!errorMessage || !errorMessage.classList.contains('error-message')) {
      errorMessage = document.createElement('p');
      errorMessage.className = 'text-error text-sm mt-1 error-message';
      errorMessage.textContent = message;
      inputElement.parentNode.insertBefore(errorMessage, inputElement.nextSibling);
    }
    
    // Focus the input
    inputElement.focus();
    
    // Remove error styling after a delay
    setTimeout(() => {
      inputElement.classList.remove('input-error');
    }, 3000);
  }
}

// Update collection via API
function updateCollection(collectionData) {
  fetch(`/collections/${collectionId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(collectionData)
  })
  .then(response => {
    if (!response.ok) {
      if (response.status === 409) {
        // Handle name conflict
        throw new Error('You already have a collection with that name');
      }
      throw new Error('Network response was not ok');
    }
    // Redirect back to the collection detail page
    window.location.href = `/collections/${collectionId}`;
  })
  .catch(error => {
    console.error('Error updating collection:', error);
    
    // Show error message to user
    const nameInput = document.getElementById('name');
    if (nameInput && error.message.includes('already have a collection')) {
      showInputError(nameInput, error.message);
    } else {
      alert('Failed to update collection: ' + error.message);
    }
  });
}

// Show delete confirmation modal
window.confirmDeleteCollection = function() {
  const modal = document.getElementById('deleteModal');
  const input = document.getElementById('confirmCollectionName');
  
  if (modal && input) {
    // Reset input field
    input.value = '';
    input.classList.remove('input-error');
    
    // Show modal
    modal.classList.add('modal-open');
  }
};

// Close delete modal
window.closeDeleteModal = function() {
  const modal = document.getElementById('deleteModal');
  if (modal) {
    modal.classList.remove('modal-open');
  }
};

// Delete collection
function deleteCollection(confirmationName) {
  fetch(`/collections/${collectionId}?name=${encodeURIComponent(confirmationName)}`, {
    method: 'DELETE'
  })
  .then(response => {
    if (!response.ok) {
      if (response.status === 409) {
        throw new Error('Confirmation name doesn\'t match');
      }
      throw new Error('Network response was not ok');
    }
    // Redirect to collections page
    window.location.href = '/collections';
  })
  .catch(error => {
    console.error('Error deleting collection:', error);
    alert('Failed to delete collection: ' + error.message);
    
    // Close the modal
    closeDeleteModal();
  });
}
