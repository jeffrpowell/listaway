require('../index')
require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const forms = document.querySelectorAll(".user-creation-form");
    const newGroupRadios = document.querySelectorAll('.newGroup');
    const existingGroupRadios = document.querySelectorAll('.existingGroup');
    const existingGroupSections = document.querySelectorAll('.existingGroupSection');
    const adminCheckboxSections = document.querySelectorAll('.adminCheckboxSection');
    const regularAdminCheckboxSections = document.querySelectorAll('.regularAdminCheckboxSection');
    const adminCheckboxes = document.querySelectorAll('.adminCheckbox');
    const errorSpans = document.querySelectorAll('.error-span');
    
    // Handle instance admin form controls
    setupInstanceAdminFormControls();

    // Take over form submission
    forms.forEach(form => form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendData(form);
    }));
    
    // Setup instance admin form controls
    function setupInstanceAdminFormControls() {
        if (newGroupRadios.length === 0 || existingGroupRadios.length === 0) {
            // Not an instance admin form, no need to set up controls
            return;
        }
        
        // Initial state - new group selected
        adminCheckboxes.forEach(checkbox => {
            checkbox.checked = true;
            checkbox.disabled = true;
        });
        
        // Set up event listeners for radio buttons
        newGroupRadios.forEach(radio => {
            radio.addEventListener('change', function() {
                if (this.checked) {
                    existingGroupSections.forEach(section => section.classList.add('hidden'));
                    adminCheckboxSections.forEach(section => section.classList.remove('hidden'));
                    regularAdminCheckboxSections.forEach(section => section.classList.add('hidden'));
                    adminCheckboxes.forEach(checkbox => {
                        checkbox.checked = true;
                        checkbox.disabled = true;
                    });
                }
            });
        });
        
        existingGroupRadios.forEach(radio => {
            radio.addEventListener('change', function() {
                if (this.checked) {
                    existingGroupSections.forEach(section => section.classList.remove('hidden'));
                    adminCheckboxSections.forEach(section => section.classList.add('hidden'));
                    regularAdminCheckboxSections.forEach(section => section.classList.remove('hidden'));
                }
            });
        });
    }

    async function sendData(form) {
        // Associate the FormData object with the form element
        const formData = new FormData(form);
    
        try {
            const response = await fetch("/admin/users/create", {
                method: "PUT",
                headers: {
                    "Accept": "text/plain",
                    "Content-Type": "application/x-www-form-urlencoded"
                },
                body: new URLSearchParams(formData).toString()
            });
            if (response.status >= 400) {
                showError(response);
            } else if (response.status === 204 || response.status === 200) {
                window.location.href = response.headers.get("Location");
            }
            else {
                showError(500);
            }
        } catch (e) {
            console.error(e);
            showError(500);
        }
    }
    
    async function showError(response) {
        text = await response.text();
        errorSpans.forEach(errorSpan => {
            if (response.status < 500) {
                errorSpan.innerText = text;
            }
            else {
                errorSpan.innerText = "Unexpected error occurred. Please try again later.";
            }
            errorSpan.classList.remove("hidden");
        });
    }
});