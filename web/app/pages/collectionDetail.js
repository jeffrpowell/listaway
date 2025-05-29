require("../navbar");
require("../collectionNav");

// Store collection ID from the page
let collectionId = null;
let availableLists = [];
let currentLists = [];

// Event handler when the page loads
document.addEventListener('DOMContentLoaded', () => {
  // Extract collection ID from URL
  const pathParts = window.location.pathname.split('/');
  collectionId = pathParts[pathParts.length - 1];
  
  // Initialize UI interactions
  initializeModals();
  
  // Get current lists in the collection
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
      id: parseInt(row.getAttribute('data-list-id')),
      displayOrder: parseInt(row.getAttribute('data-display-order'))
    };
  });
}

// Publish collection
window.publishCollection = function(id) {
  fetch(`/collections/${id}/share`, {
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
    loadAvailableLists();
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

// Load available lists that are not in the collection
function loadAvailableLists() {
  const container = document.getElementById('availableListsContainer');
  if (!container) return;
  
  container.innerHTML = '<p class="text-center py-4 text-gray-500">Loading available lists...</p>';
  
  // Get all user lists
  fetch('/lists')
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then(lists => {
      // Filter out lists already in the collection
      availableLists = lists.filter(list => 
        !currentLists.some(currentList => currentList.id === list.Id)
      );
      
      renderAvailableLists(container, availableLists);
    })
    .catch(error => {
      console.error('Error loading lists:', error);
      container.innerHTML = '<p class="text-center py-4 text-error">Error loading lists. Please try again.</p>';
    });
}

// Render available lists in the modal
function renderAvailableLists(container, lists) {
  if (lists.length === 0) {
    container.innerHTML = '<p class="text-center py-4 text-gray-500">No additional lists available to add.</p>';
    return;
  }
  
  let html = '<div class="space-y-2">';
  
  lists.forEach(list => {
    html += `
      <div class="form-control">
        <label class="label cursor-pointer justify-start">
          <input type="checkbox" class="checkbox checkbox-primary mr-2" value="${list.Id}" />
          <span class="label-text">${list.Name}</span>
        </label>
        <p class="text-xs text-gray-500 ml-6">${list.Description && list.Description.Valid ? list.Description.String : ''}</p>
      </div>
    `;
  });
  
  html += '</div>';
  container.innerHTML = html;
}

// Add selected lists to the collection
window.addSelectedLists = function() {
  const checkboxes = document.querySelectorAll('#availableListsContainer input[type="checkbox"]:checked');
  
  if (checkboxes.length === 0) {
    alert('Please select at least one list to add.');
    return;
  }
  
  const selectedListIds = Array.from(checkboxes).map(cb => parseInt(cb.value));
  let addedCount = 0;
  let nextDisplayOrder = findNextDisplayOrder();
  
  // Add each list one by one
  const addNextList = (index) => {
    if (index >= selectedListIds.length) {
      // All lists have been added
      closeAddListModal();
      window.location.reload();
      return;
    }
    
    const listId = selectedListIds[index];
    
    fetch(`/collections/${collectionId}/lists/${listId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ displayOrder: nextDisplayOrder + index })
    })
    .then(response => {
      if (!response.ok) {
        throw new Error(`Failed to add list ${listId}`);
      }
      addedCount++;
      addNextList(index + 1);
    })
    .catch(error => {
      console.error('Error adding list:', error);
      alert(`Added ${addedCount} lists, but encountered an error. Please try again.`);
      closeAddListModal();
    });
  };
  
  // Start adding lists
  addNextList(0);
};

// Find the next display order value to use
function findNextDisplayOrder() {
  if (currentLists.length === 0) return 0;
  
  const maxOrder = Math.max(...currentLists.map(list => list.displayOrder));
  return maxOrder + 1;
}

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

// Move a list up in the display order
window.moveListUp = function(listId) {
  const currentList = currentLists.find(list => list.id === listId);
  if (!currentList || currentList.displayOrder <= 0) return;
  
  const newOrder = currentList.displayOrder - 1;
  
  // Check if there's a list with this display order
  const affectedList = currentLists.find(list => list.displayOrder === newOrder);
  
  updateListOrder(listId, newOrder)
    .then(() => {
      if (affectedList) {
        // Swap positions
        return updateListOrder(affectedList.id, currentList.displayOrder);
      }
      return Promise.resolve();
    })
    .then(() => {
      window.location.reload();
    })
    .catch(error => {
      console.error('Error updating list order:', error);
      alert('Failed to update list order. Please try again.');
    });
};

// Move a list down in the display order
window.moveListDown = function(listId) {
  const currentList = currentLists.find(list => list.id === listId);
  if (!currentList) return;
  
  const newOrder = currentList.displayOrder + 1;
  
  // Check if there's a list with this display order
  const affectedList = currentLists.find(list => list.displayOrder === newOrder);
  
  updateListOrder(listId, newOrder)
    .then(() => {
      if (affectedList) {
        // Swap positions
        return updateListOrder(affectedList.id, currentList.displayOrder);
      }
      return Promise.resolve();
    })
    .then(() => {
      window.location.reload();
    })
    .catch(error => {
      console.error('Error updating list order:', error);
      alert('Failed to update list order. Please try again.');
    });
};

// Update a list's display order
function updateListOrder(listId, newOrder) {
  return fetch(`/collections/${collectionId}/lists/${listId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ displayOrder: newOrder })
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response;
  });
}
