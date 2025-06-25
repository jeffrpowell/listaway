require("../navbar")
require("../grids")
import { initMasterDetailGrid } from "../grids";

document.addEventListener('DOMContentLoaded', (event) => {
    // Initialize delete button functionality
    const deleteButtons = document.querySelectorAll('.delete-btn');
    deleteButtons.forEach(button => {
        button.addEventListener('click', async (event) => {
            event.stopPropagation(); // Prevent row expansion when clicking delete
            const listId = button.dataset.listId;
            const itemId = button.dataset.itemId;
            
            const response = await fetch('/list/' + listId + '/item/' + itemId, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            if (response.status === 204 || response.status === 200) {
                window.location.href = '/list/' + listId;
            }
        });
    });
    
    // Initialize the master/detail grid
    const gridApi = initMasterDetailGrid('.item-grid', {
        defaultSortColumn: 'priority',
        defaultSortDirection: 'asc'
    });
});