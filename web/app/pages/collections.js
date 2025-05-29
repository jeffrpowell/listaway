require("../navbar");
require("../collectionNav");

let currentCollectionId = null;
let currentCollectionName = null;

// Event handler when the page loads
document.addEventListener('DOMContentLoaded', () => {
  initializeModals();
});

// Initialize modal behavior
function initializeModals() {
  // Configure delete confirmation button
  const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
  if (confirmDeleteBtn) {
    confirmDeleteBtn.addEventListener('click', () => {
      const input = document.getElementById('confirmCollectionName');
      if (input && input.value === currentCollectionName) {
        deleteCollection(currentCollectionId, currentCollectionName);
        closeDeleteModal();
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

// Function to publish a collection (generate share code)
window.publishCollection = function(collectionId) {
  fetch(`/collections/${collectionId}/share`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    }
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.text();
  })
  .then(shareCode => {
    // Reload the page to show the share code
    window.location.reload();
  })
  .catch(error => {
    console.error('Error publishing collection:', error);
    alert('Failed to publish collection. Please try again.');
  });
};

// Function to unpublish a collection (remove share code)
window.unpublishCollection = function(collectionId) {
  fetch(`/collections/${collectionId}/share`, {
    method: 'DELETE'
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    // Reload the page to update UI
    window.location.reload();
  })
  .catch(error => {
    console.error('Error unpublishing collection:', error);
    alert('Failed to unpublish collection. Please try again.');
  });
};

// Function to display the share link modal
window.copyShareLink = function(collectionId, shareCode, sharedCollectionPath) {
  const modal = document.getElementById('shareLinkModal');
  const shareLink = document.getElementById('shareLink');
  
  if (modal && shareLink) {
    // Generate full URL for sharing
    const url = `${window.location.origin}/${sharedCollectionPath}/${shareCode}`;
    shareLink.value = url;
    
    // Display the modal
    modal.classList.add('modal-open');
  }
};

// Function to copy share link to clipboard
window.copyToClipboard = function() {
  const shareLink = document.getElementById('shareLink');
  if (shareLink) {
    shareLink.select();
    shareLink.setSelectionRange(0, 99999); // For mobile devices
    
    navigator.clipboard.writeText(shareLink.value)
      .then(() => {
        // Show success indicator
        const copyBtn = shareLink.nextElementSibling;
        if (copyBtn) {
          const originalText = copyBtn.innerHTML;
          copyBtn.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg> Copied!';
          
          setTimeout(() => {
            copyBtn.innerHTML = originalText;
          }, 2000);
        }
      })
      .catch(err => {
        console.error('Failed to copy text: ', err);
        alert('Failed to copy to clipboard. Please copy the link manually.');
      });
  }
};

// Function to close the share modal
window.closeShareModal = function() {
  const modal = document.getElementById('shareLinkModal');
  if (modal) {
    modal.classList.remove('modal-open');
  }
};

// Function to open delete confirmation modal
window.confirmDeleteCollection = function(collectionId, collectionName) {
  const modal = document.getElementById('deleteModal');
  const input = document.getElementById('confirmCollectionName');
  
  if (modal && input) {
    currentCollectionId = collectionId;
    currentCollectionName = collectionName;
    
    // Reset input field
    input.value = '';
    input.classList.remove('input-error');
    
    // Show modal
    modal.classList.add('modal-open');
  }
};

// Function to close delete modal
window.closeDeleteModal = function() {
  const modal = document.getElementById('deleteModal');
  if (modal) {
    modal.classList.remove('modal-open');
    
    // Reset state
    currentCollectionId = null;
    currentCollectionName = null;
  }
};

// Function to delete a collection
function deleteCollection(collectionId, collectionName) {
  fetch(`/collections/${collectionId}?name=${encodeURIComponent(collectionName)}`, {
    method: 'DELETE'
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    // Redirect to collections page
    window.location.href = '/collections';
  })
  .catch(error => {
    console.error('Error deleting collection:', error);
    alert('Failed to delete collection. Please try again.');
  });
}
