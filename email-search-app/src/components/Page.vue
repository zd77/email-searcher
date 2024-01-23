<template>
    <Searchbar 
        :handleEnter="handleEnter"
    />
    <div class="m-4 flex h-screen">
        <div class="w-1/2 overflow-x-auto mr-2">
            <div v-if="data && data.length > 0" class="w-auto">
                <Pagination 
                    :currentPage="currentPage"
                    :totalPages="totalPages"
                    :totalEntries="totalEntries"
                    :goToPage="goToPage"
                />
            </div>
            <Table 
                :data="data"
                :select-row="selectRow"
            />
        </div>
        <div class="w-1/2 ml-2">
            <h4 class="pl-4 font-bold">{{ subheaderContent }}</h4>
            <div class="border w-full h-full p-4">
                <p v-html="bodyMessage"></p>
            </div>
        </div>
    </div>
</template>

<script setup>
    import { ref } from 'vue'
    import Pagination from './Pagination.vue'
    import Searchbar from './SearchBar.vue'
    import Table from './Table.vue'
    const url = import.meta.env.VITE_BASE_URL
    const SEARCH_MAX_RESULT = 50
    const data = ref()
    const searchQuery = ref('')
    const subheaderContent = ref('')
    const bodyMessage = ref('...')
    const currentPage = ref(1)
    const totalEntries = ref(0)
    const totalPages = ref(1)

    const handleSearch = async ( queryString, page = 1 ) => {
        try {
            const credentials = btoa('admin:Complexpass#123')
            const response = await fetch(`${url}/api/enron_emails/_search`,{
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${credentials}`,
                },
                body: JSON.stringify({
                    "search_type": "match",
                    "query": {
                        "term": `${queryString}`,
                        "field": "_all",
                    },
                    "sort_fields": ["-@Date"],
                    "from": (page - 1) * SEARCH_MAX_RESULT,
                    "max_results": SEARCH_MAX_RESULT,
                    "_source": []
                })
            })
            const formattedResp = await response.json()
            totalEntries.value = formattedResp.hits.total.value
            totalPages.value = Math.ceil(totalEntries.value/SEARCH_MAX_RESULT)
            const extractData = formattedResp.hits.hits
            data.value = extractData.map(({_source}) => {
                const {to, from, subject, date, bodyMsg} = _source
                return {to, from, subject, date, bodyMsg}
            })
        } catch( err ) {
            console.log("Error trying to fetch: ", err)
        }
    }

    const goToPage = ( page ) => {
        if( page >= 1 && page <= totalPages.value ) {
            currentPage.value = page
            handleSearch(searchQuery.value, page )
        }
    }

    const wordHighlighter = (query, text) => {
        const regex = new RegExp(`\\b${query}\\b`, 'gi');
        const highlightedText = text.replace(regex, match => `<span class="font-bold">${match}</span>`)
        return highlightedText
    }

    const handleEnter = ( searchQ ) => {
        searchQuery.value = searchQ
        bodyMessage.value = ' '
        subheaderContent.value = ' '
        currentPage.value = 1
        handleSearch(searchQuery.value)
    }

    const selectRow = ( selectedItem ) => {
        bodyMessage.value = wordHighlighter(searchQuery.value, selectedItem.bodyMsg)
        subheaderContent.value = selectedItem.subject
    }
</script>

<style scoped>
tr:hover {
    background-color: #f0f0f0;
}
.pagination-container {
    margin-top: 16px;
  }

.pagination-button {
    margin: 0 8px;
}
</style>