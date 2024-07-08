require('../index')
require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const listNameInput = document.getElementById('listName');
    const submitButton = document.querySelector('button[type="submit"]');
    const errorSpan = document.querySelector('.error-span');

    listNameInput.addEventListener('keyup', (event) => {
        submitButton.disabled = true;
        errorSpan.classList.add('hidden');
        errorSpan.textContent = '';

        if (listNameInput.value.trim() !== '') {
            checkListName(listNameInput.value.trim());
        }
    });

    async function checkListName(name) {
        try {
            const response = await fetch(`/list/namecheck?name=${encodeURIComponent(name)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status === 200) {
                submitButton.disabled = false;
            } else if (response.status === 400) {
                submitButton.disabled = true;
                errorSpan.textContent = 'The list name is already in use.';
                errorSpan.classList.remove('hidden');
            } else {
                submitButton.disabled = true;
                errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                errorSpan.classList.remove('hidden');
            }
        } catch (error) {
            submitButton.disabled = true;
            errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
            errorSpan.classList.remove('hidden');
        }
    }
    async function sendData(form) {
        // Associate the FormData object with the form element
        const formData = new FormData(form);

        try {
            const response = await fetch("/list", {
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
            } else if (response.status === 200) {
                window.location.href = response.headers.get("Location");
            }
        } catch (e) {
            console.error(e);
            showError(500);
        }
    }

    // Take over form submission
    const forms = document.querySelectorAll(".list-form");
    forms.forEach(form => form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendData(form);
    }));

    function showError(statusCode) {
        const errorSpans = document.querySelectorAll(".error-span");
        errorSpans.forEach(errorSpan => {
            if (statusCode === 400) {
                errorSpan.innerText = "List name is already in use.";
            }
            else {
                errorSpan.innerText = "Unexpected error occurred. Please try again later.";
            }
            errorSpan.classList.remove("hidden");
        });
    }
});