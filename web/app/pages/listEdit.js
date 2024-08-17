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
    const generateShareButtons = document.querySelectorAll('.btn-generate-share');
    const shareLinks = document.querySelectorAll('.share-link');
    const copyShareLinkButtons = document.querySelectorAll('.btn-copy-share-link');
    const copyShareLinkEmptyIcons = document.querySelectorAll('.clipboard-empty');
    const copyShareLinkCheckIcons = document.querySelectorAll('.clipboard-check');
    const listItemsRedirectButtons = document.querySelectorAll('.list-items-redirect');
    const deleteListButtons = document.querySelectorAll('.list-delete');
    const deleteListConfirmationSpans = document.querySelectorAll('.list-delete-confirmation-span');
    const deleteListConfirmationInputs = document.querySelectorAll('.list-delete-confirmation');
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
            saveNameButtons.forEach(saveNameButton => saveNameButton.classList.remove("opacity-50", "cursor-not-allowed"));
            editNameErrors.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });

            if (listNameInput.value.trim() !== '') {
                await checkListName(saveNameButtons, editNameErrors, listNameInput.value.trim());
            }
            else {
                formReadyToSubmit = false;
                saveNameButtons.forEach(saveNameButton => saveNameButton.classList.add("opacity-50", "cursor-not-allowed"));
            }
        })
    );

    async function checkListName(submitButtons, errorSpans, name) {
        try {
            const response = await fetch(`/list/namecheck?name=${encodeURIComponent(name)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status === 200) {
                submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
                formReadyToSubmit = true;
            } else if (response.status === 400) {
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'The list name is already in use.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            } else {
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            }
        } catch (error) {
            submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
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
            let listId = saveNameBtn.dataset.listId;
            try {
                const response = await fetch('/list/'+listId, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(newName)
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
});