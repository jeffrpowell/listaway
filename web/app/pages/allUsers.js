require("../navbar")

document.addEventListener('DOMContentLoaded', (event) => {
    
    const adminToggles = document.querySelectorAll('.admin-toggle');
    const instanceAdminToggles = document.querySelectorAll('.instance-admin-toggle');
    const deleteUserButtons = document.querySelectorAll('.btn-delete-user');

    // Handle Group Admin toggles
    adminToggles.forEach(adminToggle => {
        adminToggle.addEventListener('click', async (event) => {
            const userId = adminToggle.dataset.userId;
            const response = await fetch('/admin/user/'+userId+'/toggleadmin', {
                method: 'POST',
            });
            if (response.status === 200) {
                adminToggle.textContent = await response.text();
            }
        });
    });
    
    // Handle Instance Admin toggles
    instanceAdminToggles.forEach(instanceAdminToggle => {
        instanceAdminToggle.addEventListener('click', async (event) => {
            const userId = instanceAdminToggle.dataset.userId;
            const response = await fetch('/admin/user/'+userId+'/toggleinstanceadmin', {
                method: 'POST',
            });
            if (response.status === 200) {
                instanceAdminToggle.textContent = await response.text();
            }
        });
    });
    
    // Handle Delete User buttons
    deleteUserButtons.forEach(deleteUserBtn => {
        deleteUserBtn.addEventListener('click', async (event) => {
            const userId = deleteUserBtn.dataset.userId;
            var firstDeleteClickDone = deleteUserBtn.dataset.deleteClicked === "true";
            if (!firstDeleteClickDone) {
                const response = await fetch('/admin/user/'+userId+'/listscount', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'text/plain'
                    },
                });
                if (response.status === 200) {
                    const listCount = await response.text();

                    // Mark the button as having been clicked once
                    deleteUserBtn.dataset.deleteClicked = "true";

                    // Find the confirmation row corresponding to this user
                    const confirmationRows = document.querySelectorAll(`.delete-confirmation-row[data-user-id="${userId}"]`);
                    // Remove the hidden class to show the confirmation row
                    confirmationRows.forEach(confirmationRow => {
                        confirmationRow.classList.remove('hidden');
                        // Find the span inside the confirmation row and update its content
                        const confirmationSpan = confirmationRow.querySelector(`.delete-confirmation-span[data-user-id="${userId}"]`);
                        confirmationSpan.textContent = `${listCount}`;
                    });
                }
            }
            else {
                const response = await fetch('/admin/user/'+userId, {
                    method: 'DELETE',
                });
                if (response.status === 204) {
                    location.reload();
                }
            }
        });
    });
});
