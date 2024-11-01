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
    const grid = document.querySelectorAll('.item-grid');
    const customTheme = themeQuartz.withParams({
        // internalContentLineHeight: "1.5"
    });

    const gridOptions = {
        theme: customTheme,
        loadThemeGoogleFonts: true,
        rowData: [],
        columnDefs: [
            { 
                field: "name", 
                headerName: "Item",
                flex: 1,
                wrapText: true,
                autoHeight: true,
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
                wrapText: true,
                autoHeight: true,
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
        let shareCode = g.dataset.shareCode;
        const gridApi = createGrid(g, gridOptions, {
            modules: [
                ClientSideRowModelModule,
            ]
        });
        fetch("/sharedlist/"+shareCode+"/items")
            .then(response => response.json())
            .then(data => gridApi.setGridOption("rowData", data));
    });
});