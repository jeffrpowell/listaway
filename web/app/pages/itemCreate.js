require('../index')
require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const forms = document.querySelectorAll(".list-form");
    const nameInputs = document.getElementsByName('name');
    const submitButtons = document.querySelectorAll('button');
    const errorSpans = document.querySelectorAll('.error-span');
    var formReadyToSubmit = nameInputs[0].value.trim() !== '';

    nameInputs.forEach(listNameInput => 
        listNameInput.addEventListener('input', (event) => {
            submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
            errorSpans.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });
            if (listNameInput.value.trim() !== '') {
                formReadyToSubmit = true;
                submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
            }
            else {
                formReadyToSubmit = false;
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
            }
        })
    );
    async function sendData(form) {
        if (!formReadyToSubmit) {
            return;
        }
        // Associate the FormData object with the form element
        const formData = new FormData(form);
        let listId = form.dataset.listId;
        let editMode = form.dataset.editMode && form.dataset.editMode === "true";
        try {
            if (editMode) {
                let itemId = form.dataset.itemId;
                await submitEditItem(listId, itemId, formData);
            }
            else {
                await submitCreateItem(listId, formData);
            }
        } catch (e) {
            console.error(e);
            showError(500);
        }
    }
    
    async function submitEditItem(listId, itemId, formData) {
        const response = await fetch("/list/"+listId+"/item/"+itemId, {
            method: "POST",
            // Set the FormData instance as the request body
            headers: {
                "Accept": "text/plain",
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: new URLSearchParams(formData).toString()
        });
        if (response.status >= 400) {
            showError(response.status);
        } else if (response.status === 204 || response.status === 200) {
            window.location.href = response.headers.get("Location");
        }
    }

    async function submitCreateItem(listId, formData) {
        const response = await fetch("/list/"+listId+"/item", {
            method: "PUT",
            // Set the FormData instance as the request body
            headers: {
                "Accept": "text/plain",
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: new URLSearchParams(formData).toString()
        });
        if (response.status >= 400) {
            showError(response.status);
        } else if (response.status === 204 || response.status === 200) {
            window.location.href = response.headers.get("Location");
        }
    }

    // Take over form submission
    forms.forEach(form => form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendData(form);
    }));

    function showError(statusCode) {
        const errorSpans = document.querySelectorAll(".error-span");
        errorSpans.forEach(errorSpan => {
            errorSpan.innerText = "Unexpected error occurred. Please try again later.";
            errorSpan.classList.remove("hidden");
        });
    }
});