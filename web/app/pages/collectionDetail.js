require("../navbar");

document.addEventListener('DOMContentLoaded', () => {
  // Initialize share links
  const shareLinks = document.querySelectorAll('.share-link');
  const listCollectionCheckboxes = document.querySelectorAll('.list-collection-checkbox');
  const copyShareLinkButtons = document.querySelectorAll('.btn-copy-share-link');
  const copyShareLinkEmptyIcons = document.querySelectorAll('.clipboard-empty');
  const copyShareLinkCheckIcons = document.querySelectorAll('.clipboard-check');
  
  shareLinks.forEach(shareLink => {
    shareLink.textContent = window.location.origin + "/" + shareLink.dataset.sharedListPath + "/" + shareLink.dataset.shareCode;
  });

  // Function to show status indicator
  const showStatus = (listId, status) => {
    const statusContainer = document.querySelector(`.request-status[data-list-id="${listId}"]`);
    if (!statusContainer) return;

    // Hide all status indicators
    statusContainer.querySelector('.success-icon').classList.add('hidden');
    statusContainer.querySelector('.loading-icon').classList.add('hidden');
    statusContainer.querySelector('.error-icon').classList.add('hidden');
    statusContainer.querySelector('.error-message').classList.add('hidden');

    // Show the requested status
    if (status === 'loading') {
      statusContainer.querySelector('.loading-icon').classList.remove('hidden');
    } else if (status === 'success') {
      statusContainer.querySelector('.success-icon').classList.remove('hidden');

      // Hide success icon after 8 seconds
      setTimeout(() => {
        const successIcon = statusContainer.querySelector('.success-icon');
        if (successIcon) {
          successIcon.classList.add('hidden');
        }
      }, 8000);
    } else if (status === 'error') {
      statusContainer.querySelector('.error-icon').classList.remove('hidden');
      statusContainer.querySelector('.error-message').classList.remove('hidden');

      // Hide error message after 8 seconds
      setTimeout(() => {
        const errorIcon = statusContainer.querySelector('.error-icon');
        const errorMessage = statusContainer.querySelector('.error-message');
        if (errorIcon && errorMessage) {
          errorIcon.classList.add('hidden');
          errorMessage.classList.add('hidden');
        }
      }, 8000);
    }
  };

  // Add event listeners to each checkbox
  listCollectionCheckboxes.forEach(checkbox => {
    checkbox.addEventListener('change', async (event) => {
      const listId = event.target.dataset.listId;
      const collectionId = event.target.dataset.collectionId;
      const hasSharecode = event.target.dataset.hasSharecode === 'true';
      const collectionHasSharecode = event.target.dataset.collectionHasSharecode === 'true';
      const isChecked = event.target.checked;

      let shouldProceed = true;

      // If the checkbox is being checked AND the list doesn't have a share code, show a confirmation
      if (collectionHasSharecode && isChecked && !hasSharecode) {
        shouldProceed = confirm(
          "Adding this list to a collection with a share code will generate a share code for this list as well. " +
          "This means the list will be publicly accessible via its share link. " +
          "Do you want to continue?"
        );
        
        if (!shouldProceed) {
          // User cancelled, revert checkbox state
          event.target.checked = false;
          return;
        }
      }

      // Show loading indicator
      showStatus(listId, 'loading');

      try {
        let response;
        if (isChecked) {
          // Add list to collection
          response = await fetch(`/collections/${collectionId}/lists/${listId}`, {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json'
            }
          });
        } else {
          // Remove list from collection
          response = await fetch(`/collections/${collectionId}/lists/${listId}`, {
            method: 'DELETE',
            headers: {
              'Content-Type': 'application/json'
            }
          });
        }

        // Check if response was successful
        if (response.ok) {
          showStatus(listId, 'success');
          
          // If we just added a list to a collection and generated its share code, update the data attribute
          if (isChecked && !hasSharecode) {
            event.target.dataset.hasSharecode = 'true';
            
            // Refresh the page after a brief delay to show the new share code
            setTimeout(() => {
              window.location.reload();
            }, 1500);
          }
        } else {
          throw new Error(`Request failed with status ${response.status}`);
        }
      } catch (error) {
        console.error('Error updating collection:', error);
        // Revert checkbox state on error
        event.target.checked = !isChecked;
        showStatus(listId, 'error');
      }
    });
  });

  copyShareLinkButtons.forEach(copyShareLinkBtn => {
    copyShareLinkBtn.addEventListener('click', async (event) => {
      var result = writeClipboardText(window.location.origin + "/" + copyShareLinkBtn.dataset.sharedListPath + "/" + copyShareLinkBtn.dataset.shareCode);
      if (result) {
        copyShareLinkEmptyIcons.forEach(icon => {
          if (icon.getAttribute('data-share-code') === copyShareLinkBtn.dataset.shareCode) {
            icon.classList.add("hidden")
          }
        });
        copyShareLinkCheckIcons.forEach(icon => {
          if (icon.getAttribute('data-share-code') === copyShareLinkBtn.dataset.shareCode) {
            icon.classList.remove("hidden")
          }
        });
      }
    });
  });

  async function writeClipboardText(text) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch (error) {
      console.error(error.message);
      return false;
    }
  }
});