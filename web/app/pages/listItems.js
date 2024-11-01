require('../index')
require('../navbar')
import { themeQuartz } from '@ag-grid-community/theming';
import { ClientSideRowModelModule } from "@ag-grid-community/client-side-row-model";
import { createGrid } from "@ag-grid-community/core";

const itemNameRenderer = function(params) {
    if (params.data.url.Valid) {
        return `<a href="${params.data.url.String}" class="underline" target="_blank" rel="nofollow noreferrer">${params.value}</a>`
    } else {
        return params.value;
    }
}

document.addEventListener('DOMContentLoaded', (event) => {
    const deleteItemButtons = document.querySelectorAll('.btn-delete-item');
    const grid = document.querySelectorAll('.item-grid');

    const gridOptions = {
        theme: themeQuartz,
        loadThemeGoogleFonts: true,
        rowData: [],
        columnDefs: [
            { 
                field: "name", 
                headerName: "Item",
                flex: 1,
                cellRenderer: itemNameRenderer
            },
            { 
                field: "priority", 
                headerName: "Priority", 
                type: 'numericColumn',
                width: 125,
                sort: 'asc',
                valueGetter: p => p.data.priority.Valid ? p.data.priority.Int64 : null
            },
            { 
                field: "notes", 
                headerName: "Notes", 
                flex: 2,
                valueGetter: p => p.data.notes.Valid ? p.data.notes.String : ""
            },
        ],
        domLayout: 'autoHeight',
        getRowId: params => String(params.data.id),
        autoSizeStrategy: {
            type: "fitGridWidth",
        }
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