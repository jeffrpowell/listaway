require("../navbar");

// Store collection ID from the page
let collectionId = null;
let currentLists = [];

// Event handler when the page loads
document.addEventListener('DOMContentLoaded', () => {
  // Extract collection ID from URL
  const pathParts = window.location.pathname.split('/');
  collectionId = pathParts[pathParts.length - 1];
  
  // Initialize UI interactions
  initializeModals();
  
  // Get current lists in the collection (already loaded from server)
  getCurrentLists();
});

// Initialize modal behavior
function initializeModals() {
  // Setup handlers for modals if they exist
  setupShareModal();
}

// Setup share modal
function setupShareModal() {
  const modal = document.getElementById('shareLinkModal');
  if (!modal) return;
  
  // Close when clicking outside the modal box
  modal.addEventListener('click', function(e) {
    if (e.target === this) {
      closeShareModal();
    }
  });
}

// Get current lists in the collection
function getCurrentLists() {
  const listRows = document.querySelectorAll('#listTable tr[data-list-id]');
  currentLists = Array.from(listRows).map(row => {
    return {
      id: parseInt(row.getAttribute('data-list-id'))
    };
  });
}

// Publish collection
window.publishCollection = function(id) {
  // Create and submit a form to publish the collection
  const form = document.createElement('form');
  form.method = 'POST';
  form.action = `/collections/${id}/share`;
  
  document.body.appendChild(form);
  form.submit();
};

// Share link handling
window.copyShareLink = function(shareCode) {
  const modal = document.getElementById('shareLinkModal');
  const shareLink = document.getElementById('shareLink');
  
  if (modal && shareLink) {
    // Generate full URL for sharing
    const url = `${window.location.origin}/sharedcollection/${shareCode}`;
    shareLink.value = url;
    
    // Display the modal
    modal.classList.add('modal-open');
  }
};

// Copy to clipboard
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

// Close share modal
window.closeShareModal = function() {
  const modal = document.getElementById('shareLinkModal');
  if (modal) {
    modal.classList.remove('modal-open');
  }
};

// Open add list modal
window.openAddListModal = function() {
  const modal = document.getElementById('addListModal');
  if (modal) {
    // Available lists are now loaded from server
    modal.classList.add('modal-open');
  }
};

// Close add list modal
window.closeAddListModal = function() {
  const modal = document.getElementById('addListModal');
  if (modal) {
    modal.classList.remove('modal-open');
  }
};

// Add selected lists to the collection
window.addSelectedLists = function() {
  const checkboxes = document.querySelectorAll('#availableListsContainer input[type="checkbox"]:checked');
  
  if (checkboxes.length === 0) {
    alert('Please select at least one list to add.');
    return;
  }
  
  const selectedListIds = Array.from(checkboxes).map(cb => parseInt(cb.value));
  const form = document.createElement('form');
  form.method = 'POST';
  form.action = `/collections/${collectionId}/lists`;
  
  // Add inputs for each selected list ID
  selectedListIds.forEach(listId => {
    const input = document.createElement('input');
    input.type = 'hidden';
    input.name = 'listIds';
    input.value = listId;
    form.appendChild(input);
  });
  
  // Submit the form
  document.body.appendChild(form);
  form.submit();
};

// Remove a list from the collection
window.removeListFromCollection = function(listId) {
  if (!confirm('Are you sure you want to remove this list from the collection?')) {
    return;
  }
  
  fetch(`/collections/${collectionId}/lists/${listId}`, {
    method: 'DELETE'
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    // Reload the page to update the list
    window.location.reload();
  })
  .catch(error => {
    console.error('Error removing list:', error);
    alert('Failed to remove list from collection. Please try again.');
  });
};