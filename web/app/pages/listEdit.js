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
    var formReadyToSubmit = false;

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
        listNameInput.addEventListener('input', (event) => {
            saveNameButtons.forEach(saveNameButton => saveNameButton.classList.remove("opacity-50", "cursor-not-allowed"));
            editNameErrors.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });

            if (listNameInput.value.trim() !== '') {
                checkListName(saveNameButtons, editNameErrors, listNameInput.value.trim());
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
    
    // listNameInputs.forEach(listNameInput => 
    //     listNameInput.addEventListener('input', (event) => {
    //         submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
    //         errorSpans.forEach(errorSpan => {
    //             errorSpan.classList.add('hidden');
    //             errorSpan.textContent = '';
    //         });

    //         if (listNameInput.value.trim() !== '') {
    //             checkListName(submitButtons, errorSpans, listNameInput.value.trim());
    //         }
    //     })
    // );

    // async function checkListName(submitButtons, errorSpans, name) {
    //     try {
    //         const response = await fetch(`/list/namecheck?name=${encodeURIComponent(name)}`, {
    //             method: 'GET',
    //             headers: {
    //                 'Content-Type': 'application/json'
    //             }
    //         });

    //         if (response.status === 200) {
    //             submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
    //             formReadyToSubmit = true;
    //         } else if (response.status === 400) {
    //             submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
    //             errorSpans.forEach(errorSpan => {
    //                 errorSpan.textContent = 'The list name is already in use.';
    //                 errorSpan.classList.remove('hidden');
    //             });
    //             formReadyToSubmit = false;
    //         } else {
    //             submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
    //             errorSpans.forEach(errorSpan => {
    //                 errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
    //                 errorSpan.classList.remove('hidden');
    //             });
    //             formReadyToSubmit = false;
    //         }
    //     } catch (error) {
    //         submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
    //         errorSpans.forEach(errorSpan => {
    //             errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
    //             errorSpan.classList.remove('hidden');
    //         });
    //         formReadyToSubmit = false;
    //     }
    // }
    // async function sendData(form) {
    //     if (!formReadyToSubmit) {
    //         return;
    //     }
    //     // Associate the FormData object with the form element
    //     const formData = new FormData(form);

    //     try {
    //         const response = await fetch("/list", {
    //             method: "PUT",
    //             // Set the FormData instance as the request body
    //             headers: {
    //                 "Accept": "text/plain",
    //                 "Content-Type": "application/x-www-form-urlencoded"
    //             },
    //             body: new URLSearchParams(formData).toString()
    //         });
    //         if (response.status >= 400) {
    //             showError(response.status);
    //         } else if (response.status === 200) {
    //             window.location.href = response.headers.get("Location");
    //         }
    //     } catch (e) {
    //         console.error(e);
    //         showError(500);
    //     }
    // }

    // // Take over form submission
    // const forms = document.querySelectorAll(".list-form");
    // forms.forEach(form => form.addEventListener("submit", (event) => {
    //     event.preventDefault();
    //     sendData(form);
    // }));

    // function showError(statusCode) {
    //     const errorSpans = document.querySelectorAll(".error-span");
    //     errorSpans.forEach(errorSpan => {
    //         if (statusCode === 400) {
    //             errorSpan.innerText = "List name is already in use.";
    //         }
    //         else {
    //             errorSpan.innerText = "Unexpected error occurred. Please try again later.";
    //         }
    //         errorSpan.classList.remove("hidden");
    //     });
    // }
});