require('../index')
require('../navbar')
import "@ag-grid-community/styles/ag-grid.min.css";
import "@ag-grid-community/styles/ag-theme-quartz.min.css";
import { ClientSideRowModelModule } from "@ag-grid-community/client-side-row-model";
import { createGrid } from "@ag-grid-community/core";

document.addEventListener('DOMContentLoaded', (event) => {
    const deleteItemButtons = document.querySelectorAll('.btn-delete-item');
    const grid = document.querySelectorAll('.item-grid');

    const gridOptions = {
        rowData: [],
        columnDefs: [
            { 
                field: "name", 
                headerName: "Item"
            },
            { 
                field: "url", 
                headerName: "URL", 
                valueGetter: p => p.data.Valid ? p.data.String : ""
            },
            { 
                field: "priority", 
                headerName: "Priority", 
                type: 'numericColumn',
                valueGetter: p => p.data.Valid ? p.data.Int64 : null
            },
            { 
                field: "notes", 
                headerName: "Notes", 
                valueGetter: p => p.data.Valid ? p.data.String : ""
            },
        ],
        getRowId: params => String(params.data.id)
    };
    grid.forEach(g => {
        const gridApi = createGrid(g, gridOptions, {
            modules: [
                ClientSideRowModelModule,
            ]
        });
        let listId = g.dataset.listId;
        fetch("/list/"+listId+"/items")
            .then(response => response.json())
            .then(data => gridApi.setGridOption("rowData", data));
    });

    deleteItemButtons.forEach(deleteItem => {
        deleteItem.addEventListener('click', async (event) => {
            let listId = deleteItem.dataset.listId;
            let itemId = deleteItem.dataset.itemId;
            const response = await fetch('/list/' + listId + '/item/' + itemId, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (response.status === 204 || response.status === 200) {
                window.location.href = '/list/' + listId;
            }
        });
    });

});