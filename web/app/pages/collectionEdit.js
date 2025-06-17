require('../index')
require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const collectionNameHeaders = document.querySelectorAll('.collection-name-header');
    const collectionNameInputs = document.querySelectorAll('.collection-name-input');
    const editNameActions = document.querySelectorAll('.edit-name-edit-actions');
    const saveNameButtons = document.querySelectorAll('.btn-save-collection-name');
    const editNameSpinners = document.querySelectorAll('.edit-name-spinner');
    const cancelNameButtons = document.querySelectorAll('.btn-cancel-collection-name');
    const readNameActions = document.querySelectorAll('.edit-name-read-actions');
    const editNameButtons = document.querySelectorAll('.btn-edit-collection-name');
    const editNameErrors = document.querySelectorAll('.edit-name-error');
    const collectionDescriptionInputs = document.querySelectorAll('.collection-description-input');
    const descriptionSavedIcons = document.querySelectorAll('.description-saved-icon');
    const descriptionErrors = document.querySelectorAll('.description-error');
    const generateShareButtons = document.querySelectorAll('.btn-generate-share');
    const shareLinks = document.querySelectorAll('.share-link');
    const copyShareLinkButtons = document.querySelectorAll('.btn-copy-share-link');
    const copyShareLinkEmptyIcons = document.querySelectorAll('.clipboard-empty');
    const copyShareLinkCheckIcons = document.querySelectorAll('.clipboard-check');
    const unpublishShareButtons = document.querySelectorAll('.btn-unpublish-share');
    const collectionItemsRedirectButtons = document.querySelectorAll('.collection-items-redirect');
    const deleteCollectionButtons = document.querySelectorAll('.collection-delete');
    const deleteCollectionConfirmationSpans = document.querySelectorAll('.collection-delete-confirmation-span');
    const deleteCollectionConfirmationInputs = document.querySelectorAll('.collection-delete-confirmation');
    var formReadyToSubmit = false;
    var firstDeleteClickDone = false;

    editNameButtons.forEach(editNameBtn => {
        editNameBtn.addEventListener('click', (event) => {
            collectionNameHeaders.forEach(el => el.classList.add('hidden'));
            collectionNameInputs.forEach(el => el.classList.remove('hidden'));
            readNameActions.forEach(el => el.classList.add('hidden'));
            editNameActions.forEach(el => el.classList.remove('hidden'));
        });
    });
    
    cancelNameButtons.forEach(cancelNameBtn => {
        cancelNameBtn.addEventListener('click', (event) => {
            collectionNameHeaders.forEach(el => el.classList.remove('hidden'));
            collectionNameInputs.forEach(el => el.classList.add('hidden'));
            editNameActions.forEach(el => el.classList.add('hidden'));
            editNameSpinners.forEach(el => el.classList.add('hidden'));
            readNameActions.forEach(el => el.classList.remove('hidden'));
        });
    });

    collectionNameInputs.forEach(collectionNameInput => 
        collectionNameInput.addEventListener('input', async (event) => {
            editNameErrors.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });
            saveNameButtons.forEach(saveNameButton => saveNameButton.classList.add("opacity-50", "cursor-not-allowed"));
            if (collectionNameInput.value.trim() !== '') {
                debouncedCheckCollectionName(editNameErrors, collectionNameInput.value.trim());
            }
            else {
                formReadyToSubmit = false;
            }
        })
    );

    var debouncedCheckCollectionName = debounce(checkCollectionName, 500);

    async function checkCollectionName(errorSpans, name) {
        try {
            const response = await fetch(`/collections/namecheck?name=${encodeURIComponent(name)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status === 200) {
                saveNameButtons.forEach(saveNameButton => saveNameButton.classList.remove("opacity-50", "cursor-not-allowed"));
                formReadyToSubmit = true;
                
                // Use one-time event listener to avoid multiple bindings
                saveNameButtons.forEach(saveNameButton => {
                    // Remove previous event listeners
                    const newBtn = saveNameButton.cloneNode(true);
                    saveNameButton.parentNode.replaceChild(newBtn, saveNameButton);
                    
                    newBtn.addEventListener('click', async (event) => {
                        if (!formReadyToSubmit) {
                            return;
                        }
                        
                        const collectionId = saveNameButton.dataset.collectionId;
                        let collectionName;
                        collectionNameInputs.forEach(collectionNameInput => {
                            collectionName = collectionNameInput.value.trim();
                        });

                        let description;
                        for (const collectionDescriptionInput of collectionDescriptionInputs) {
                            description = collectionDescriptionInput.value;
                            break;
                        }
                        
                        if (collectionName) {
                            // Hide the form error message if it exists
                            editNameErrors.forEach(errorSpan => {
                                errorSpan.classList.add('hidden');
                                errorSpan.textContent = '';
                            });
                            
                            editNameSpinners.forEach(el => el.classList.remove('hidden'));
                            
                            try {
                                const response = await fetch(`/collections/${collectionId}`, {
                                    method: 'PUT',
                                    headers: {
                                        'Content-Type': 'application/json'
                                    },
                                    body: JSON.stringify({
                                        name: collectionName,
                                        description: description
                                    })
                                });
                                
                                if (response.status === 204) {
                                    collectionNameHeaders.forEach(el => {
                                        el.textContent = collectionName;
                                        el.classList.remove('hidden');
                                    });
                                    collectionNameInputs.forEach(el => el.classList.add('hidden'));
                                    editNameActions.forEach(el => el.classList.add('hidden'));
                                    readNameActions.forEach(el => el.classList.remove('hidden'));
                                } else {
                                    editNameErrors.forEach(errorSpan => {
                                        errorSpan.classList.remove('hidden');
                                    });
                                }
                                formReadyToSubmit = false;
                            } catch (error) {
                                console.error('Error saving collection name:', error);
                                editNameErrors.forEach(errorSpan => {
                                    errorSpan.classList.remove('hidden');
                                });
                            } finally {
                                editNameSpinners.forEach(el => el.classList.add('hidden'));
                            }
                        }
                    })
                });
            } else {
                formReadyToSubmit = false;
                
                // Show error message for name conflict
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'A collection with this name already exists.';
                    errorSpan.classList.remove('hidden');
                });
            }
        } catch (error) {
            console.error('Error checking collection name:', error);
            formReadyToSubmit = false;
        }
    }

    // Auto-save description on input changes with debounce
    collectionDescriptionInputs.forEach(descriptionInput => {
        let saveTimeout;
        descriptionInput.addEventListener('input', (event) => {
            clearTimeout(saveTimeout);
            saveTimeout = setTimeout(() => {
                saveDescription(descriptionInput.dataset.collectionId, descriptionInput.value);
            }, 1000);
        });
    });

    async function saveDescription(collectionId, description) {
        try {
            let collectionName;
            for (const collectionNameHeader of collectionNameHeaders) {
                collectionName = collectionNameHeader.textContent;
                break;
            }
            const response = await fetch(`/collections/${collectionId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: collectionName,
                    description: description
                })
            });

            if (response.status === 204) {
                // Show saved confirmation
                descriptionSavedIcons.forEach(icon => {
                    icon.classList.remove('hidden');
                    setTimeout(() => {
                        icon.classList.add('hidden');
                    }, 2000);
                });
            } else {
                // Show error message
                descriptionErrors.forEach(errorElement => {
                    errorElement.classList.remove('hidden');
                    setTimeout(() => {
                        errorElement.classList.add('hidden');
                    }, 3000);
                });
            }
        } catch (error) {
            console.error('Error saving description:', error);
            // Show error message
            descriptionErrors.forEach(errorElement => {
                errorElement.classList.remove('hidden');
                setTimeout(() => {
                    errorElement.classList.add('hidden');
                }, 3000);
            });
        }
    }

    // Handle share links
    generateShareButtons.forEach(generateShareBtn => {
        generateShareBtn.addEventListener('click', async (event) => {
            let collectionId = generateShareBtn.dataset.collectionId;
            const response = await fetch('/collections/'+collectionId+'/share', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (response.status === 200) {
                window.location.reload();
            }
        });
    });
    
    shareLinks.forEach(shareLink => {
        shareLink.textContent = window.location.origin + "/" + shareLink.dataset.sharedCollectionPath + "/" + shareLink.dataset.shareCode;
    });

    copyShareLinkButtons.forEach(copyShareLinkBtn => {
        copyShareLinkBtn.addEventListener('click', async (event) => {
            var result = writeClipboardText(window.location.origin + "/" + copyShareLinkBtn.dataset.sharedCollectionPath + "/" + copyShareLinkBtn.dataset.shareCode);
            if (result) {
                copyShareLinkEmptyIcons.forEach(icon => icon.classList.add("hidden"));
                copyShareLinkCheckIcons.forEach(icon => icon.classList.remove("hidden"));
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

    unpublishShareButtons.forEach(unpublishShareBtn => {
        unpublishShareBtn.addEventListener('click', async (event) => {
            let collectionId = unpublishShareBtn.dataset.collectionId;
            const response = await fetch('/collections/'+collectionId+'/share', {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (response.status === 204 || response.status === 200) {
                window.location.reload();
            }
        });
    });
    
    collectionItemsRedirectButtons.forEach(collectionItemsRedirectBtn => {
        collectionItemsRedirectBtn.addEventListener('click', async (event) => {
            let collectionId = collectionItemsRedirectBtn.dataset.collectionId;
            window.location.href = "/collections/"+collectionId;
        });
    });

    deleteCollectionButtons.forEach(deleteCollectionBtn => {
        deleteCollectionBtn.addEventListener('click', async (event) => {
            if (!firstDeleteClickDone) {
                deleteCollectionConfirmationSpans.forEach(el => el.classList.remove('hidden'));
                firstDeleteClickDone = true;
            }
            else {
                let trueName;
                for (const collectionNameHeader of collectionNameHeaders) {
                    trueName = collectionNameHeader.textContent;
                    break;
                }
                let confirmName;
                for (const input of deleteCollectionConfirmationInputs) {
                    if (input.value === trueName) {
                        confirmName = input.value;
                        break;
                    }
                }
                if (confirmName === undefined || confirmName === null) {
                    return;
                }
                let collectionId = deleteCollectionBtn.dataset.collectionId;
                // Fix: Proper format for passing the confirmation name
                const response = await fetch('/collections/'+collectionId+'?name='+encodeURIComponent(confirmName), {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                if (response.status === 204) {
                    window.location.href = response.headers.get("Location") || '/list';
                }
            }
        });
    });

    function debounce(func, delay) {
        let timeoutId;
        return function(...args) {
            if (timeoutId) {
                clearTimeout(timeoutId);
            }
            timeoutId = setTimeout(() => {
                func.apply(this, args);
            }, delay);
        };
    }
});
