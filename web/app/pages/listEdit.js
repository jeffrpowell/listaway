require('../index')
require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const listNameHeaders = document.querySelectorAll('.list-name-header');
    const listNameInputs = document.querySelectorAll('.list-name-input');
    const editNameActions = document.querySelectorAll('.edit-name-edit-actions');
    const saveNameButtons = document.querySelectorAll('.btn-save-list-name');
    const editNameSpinners = document.querySelectorAll('.edit-name-spinner');
    const cancelNameButtons = document.querySelectorAll('.btn-cancel-list-name');
    const readNameActions = document.querySelectorAll('.edit-name-read-actions');
    const editNameButtons = document.querySelectorAll('.btn-edit-list-name');
    const editNameErrors = document.querySelectorAll('.edit-name-error');
    const listDescriptionInputs = document.querySelectorAll('.list-description-input');
    const descriptionSavedIcons = document.querySelectorAll('.description-saved-icon');
    const descriptionErrors = document.querySelectorAll('.description-error');
    const generateShareButtons = document.querySelectorAll('.btn-generate-share');
    const shareLinks = document.querySelectorAll('.share-link');
    const copyShareLinkButtons = document.querySelectorAll('.btn-copy-share-link');
    const copyShareLinkEmptyIcons = document.querySelectorAll('.clipboard-empty');
    const copyShareLinkCheckIcons = document.querySelectorAll('.clipboard-check');
    const unpublishShareButtons = document.querySelectorAll('.btn-unpublish-share');
    const listItemsRedirectButtons = document.querySelectorAll('.list-items-redirect');
    const deleteListButtons = document.querySelectorAll('.list-delete');
    const deleteListConfirmationSpans = document.querySelectorAll('.list-delete-confirmation-span');
    const deleteListConfirmationInputs = document.querySelectorAll('.list-delete-confirmation');
    const shareWithGroupCheckboxes = document.querySelectorAll('.checkbox-share-with-group');
    const groupCanEditCheckboxes = document.querySelectorAll('.checkbox-group-can-edit');
    const groupSharingStatus = document.querySelectorAll('.group-sharing-status');
    const groupSharingError = document.querySelectorAll('.group-sharing-error');
    var formReadyToSubmit = false;
    var firstDeleteClickDone = false;

    editNameButtons.forEach(editNameBtn => {
        editNameBtn.addEventListener('click', (event) => {
            listNameHeaders.forEach(el => el.classList.add('hidden'));
            listNameInputs.forEach(el => el.classList.remove('hidden'));
            readNameActions.forEach(el => el.classList.add('hidden'));
            editNameActions.forEach(el => el.classList.remove('hidden'));
        });
    });
    
    cancelNameButtons.forEach(cancelNameBtn => {
        cancelNameBtn.addEventListener('click', (event) => {
            listNameHeaders.forEach(el => el.classList.remove('hidden'));
            listNameInputs.forEach(el => el.classList.add('hidden'));
            editNameActions.forEach(el => el.classList.add('hidden'));
            editNameSpinners.forEach(el => el.classList.add('hidden'));
            readNameActions.forEach(el => el.classList.remove('hidden'));
        });
    });

    listNameInputs.forEach(listNameInput => 
        listNameInput.addEventListener('input', async (event) => {
            editNameErrors.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });
            saveNameButtons.forEach(saveNameButton => saveNameButton.classList.add("opacity-50", "cursor-not-allowed"));
            if (listNameInput.value.trim() !== '') {
                debouncedCheckListName(editNameErrors, listNameInput.value.trim());
            }
            else {
                formReadyToSubmit = false;
            }
        })
    );

    var debouncedCheckListName = debounce(checkListName, 500);
    async function checkListName(errorSpans, name) {
        try {
            const response = await fetch(`/list/namecheck?name=${encodeURIComponent(name)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status === 200) {
                saveNameButtons.forEach(saveNameButton => saveNameButton.classList.remove("opacity-50", "cursor-not-allowed"));
                formReadyToSubmit = true;
            } else if (response.status === 400) {
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'The list name is already in use.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            } else {
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            }
        } catch (error) {
            errorSpans.forEach(errorSpan => {
                errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                errorSpan.classList.remove('hidden');
            });
            formReadyToSubmit = false;
        }
    }
    
    saveNameButtons.forEach(saveNameBtn => {
        saveNameBtn.addEventListener('click', async (event) => {
            if (!formReadyToSubmit) {
                return;
            }
            saveNameButtons.forEach(el => el.classList.add('hidden'));
            editNameSpinners.forEach(el => el.classList.remove('hidden'));
            let oldName;
            for (const listNameHeader of listNameHeaders) {
                oldName = listNameHeader.textContent;
                break;
            }
            let newName;
            for (const input of listNameInputs) {
                if (input.value !== oldName) {
                    newName = input.value;
                    break;
                }
            }
            let description;
            for (const listDescriptionInput of listDescriptionInputs) {
                description = listDescriptionInput.value;
                break;
            }
            
            let shareWithGroup = false;
            for (const checkbox of shareWithGroupCheckboxes) {
                shareWithGroup = checkbox.checked;
                break;
            }
            
            let groupCanEdit = false;
            for (const checkbox of groupCanEditCheckboxes) {
                groupCanEdit = checkbox.checked;
                break;
            }
            
            let listId = saveNameBtn.dataset.listId;
            try {
                const response = await fetch('/list/'+listId, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: newName, 
                        description: description,
                        shareWithGroup: shareWithGroup,
                        groupCanEdit: groupCanEdit
                    })
                });
                
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                
                // Handle success
                listNameHeaders.forEach(el => el.textContent = newName);
                editNameSpinners.forEach(el => el.classList.add('hidden'));
                saveNameButtons.forEach(el => el.classList.remove('hidden'));
                editNameActions.forEach(el => el.classList.add('hidden'));
                readNameActions.forEach(el => el.classList.remove('hidden'));
                listNameHeaders.forEach(el => el.classList.remove('hidden'));
                listNameInputs.forEach(el => el.classList.add('hidden'));
            } catch (error) {
                // Handle error
                editNameSpinners.forEach(el => el.classList.add('hidden'));
                saveNameButtons.forEach(el => el.classList.remove('hidden'));
                saveNameButtons.forEach(el => el.classList.remove('hidden'));
                editNameErrors.forEach(el => el.classList.remove('hidden'));
            }
        });
    });

    listDescriptionInputs.forEach(listDescriptionInput => {
        listDescriptionInput.addEventListener('input', async (event) => {
            descriptionSavedIcons.forEach(descriptionSavedIcon => descriptionSavedIcon.classList.add('hidden'));
            descriptionErrors.forEach(descriptionError => descriptionError.classList.add('hidden'));
            debouncedSaveDescription(listDescriptionInput.dataset.listId, listDescriptionInput.value.trim());
        });
        
        // Save immediately when user leaves the field to ensure changes aren't lost
        listDescriptionInput.addEventListener('blur', async (event) => {
            // Cancel any pending debounced save and save immediately
            if (debouncedSaveDescription.timeoutId) {
                clearTimeout(debouncedSaveDescription.timeoutId);
            }
            await saveDescription(listDescriptionInput.dataset.listId, listDescriptionInput.value.trim());
        });
    });

    var debouncedSaveDescription = debounce(saveDescription, 500);
    async function saveDescription(listId, description) {
        let listName;
        for (const listNameHeader of listNameHeaders) {
            listName = listNameHeader.textContent;
            break;
        }
        
        let shareWithGroup = false;
        for (const checkbox of shareWithGroupCheckboxes) {
            shareWithGroup = checkbox.checked;
            break;
        }
        
        let groupCanEdit = false;
        for (const checkbox of groupCanEditCheckboxes) {
            groupCanEdit = checkbox.checked;
            break;
        }
        
        try {
            const response = await fetch('/list/'+listId, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: listName, 
                    description: description,
                    shareWithGroup: shareWithGroup,
                    groupCanEdit: groupCanEdit
                })
            });
            
            if (response.status !== 204 && response.status !== 200) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            descriptionSavedIcons.forEach(descriptionSavedIcon => descriptionSavedIcon.classList.remove('hidden'));
            setTimeout(() => descriptionSavedIcons.forEach(descriptionSavedIcon => descriptionSavedIcon.classList.add('hidden')), 5000);
            listDescriptionInputs.forEach(listDescriptionInput => listDescriptionInput.value = description);
        }
        catch (error) {
            descriptionErrors.forEach(descriptionError => descriptionError.classList.remove('hidden'));
        }
    }
    
    generateShareButtons.forEach(generateShareBtn => {
        generateShareBtn.addEventListener('click', async (event) => {
            let listId = generateShareBtn.dataset.listId;
            const response = await fetch('/list/'+listId+'/share', {
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
        shareLink.textContent = window.location.origin + "/" + shareLink.dataset.sharedListPath + "/" + shareLink.dataset.shareCode;
    });

    copyShareLinkButtons.forEach(copyShareLinkBtn => {
        copyShareLinkBtn.addEventListener('click', async (event) => {
            var result = writeClipboardText(window.location.origin + "/" + copyShareLinkBtn.dataset.sharedListPath + "/" + copyShareLinkBtn.dataset.shareCode);
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
            let listId = unpublishShareBtn.dataset.listId;
            const response = await fetch('/list/'+listId+'/share', {
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
    
    listItemsRedirectButtons.forEach(listItemsRedirectBtn => {
        listItemsRedirectBtn.addEventListener('click', async (event) => {
            let listId = listItemsRedirectBtn.dataset.listId;
            window.location.href = "/list/"+listId;
        });
    });

    deleteListButtons.forEach(deleteListBtn => {
        deleteListBtn.addEventListener('click', async (event) => {
            if (!firstDeleteClickDone) {
                deleteListConfirmationSpans.forEach(el => el.classList.remove('hidden'));
                firstDeleteClickDone = true;
            }
            else {
                let trueName;
                for (const listNameHeader of listNameHeaders) {
                    trueName = listNameHeader.textContent;
                    break;
                }
                let confirmName;
                for (const input of deleteListConfirmationInputs) {
                    if (input.value === trueName) {
                        confirmName = input.value;
                        break;
                    }
                }
                if (confirmName === undefined || confirmName === null) {
                    return;
                }
                let listId = deleteListBtn.dataset.listId;
                const response = await fetch('/list/'+listId, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(confirmName)
                });
                if (response.status === 200) {
                    window.location.href = response.headers.get("Location");
                }
            }
        });
    });

    // Handle Enter key in delete confirmation input
    deleteListConfirmationInputs.forEach(input => {
        input.addEventListener('keydown', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                deleteListButtons.forEach(btn => btn.click());
            }
        });
    });

    // Group sharing checkboxes
    shareWithGroupCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', async (event) => {
            const isChecked = checkbox.checked;
            const listId = checkbox.dataset.listId;
            
            // Enable/disable the group can edit checkbox based on share with group
            groupCanEditCheckboxes.forEach(editCheckbox => {
                if (isChecked) {
                    editCheckbox.disabled = false;
                } else {
                    editCheckbox.disabled = true;
                    editCheckbox.checked = false;
                }
            });
            
            await saveGroupSharingSettings(listId);
        });
    });
    
    groupCanEditCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', async (event) => {
            const listId = checkbox.dataset.listId;
            await saveGroupSharingSettings(listId);
        });
    });
    
    async function saveGroupSharingSettings(listId) {
        groupSharingStatus.forEach(el => el.classList.add('hidden'));
        groupSharingError.forEach(el => el.classList.add('hidden'));
        
        let listName;
        for (const listNameHeader of listNameHeaders) {
            listName = listNameHeader.textContent;
            break;
        }
        
        let description;
        for (const listDescriptionInput of listDescriptionInputs) {
            description = listDescriptionInput.value;
            break;
        }
        
        let shareWithGroup = false;
        for (const checkbox of shareWithGroupCheckboxes) {
            shareWithGroup = checkbox.checked;
            break;
        }
        
        let groupCanEdit = false;
        for (const checkbox of groupCanEditCheckboxes) {
            groupCanEdit = checkbox.checked;
            break;
        }
        
        try {
            const response = await fetch('/list/' + listId, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: listName,
                    description: description,
                    shareWithGroup: shareWithGroup,
                    groupCanEdit: groupCanEdit
                })
            });
            
            if (response.status !== 204 && response.status !== 200) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            groupSharingStatus.forEach(el => el.classList.remove('hidden'));
            setTimeout(() => groupSharingStatus.forEach(el => el.classList.add('hidden')), 3000);
        } catch (error) {
            groupSharingError.forEach(el => el.classList.remove('hidden'));
        }
    }

    function debounce(func, delay) {
        let timeoutId;
        const debouncedFunc = function(...args) {
            if (timeoutId) {
                clearTimeout(timeoutId);
            }
            timeoutId = setTimeout(() => {
                func.apply(this, args);
            }, delay);
        };
        // Expose timeoutId for external access
        Object.defineProperty(debouncedFunc, 'timeoutId', {
            get: () => timeoutId,
            set: (value) => { timeoutId = value; }
        });
        return debouncedFunc;
    }
});