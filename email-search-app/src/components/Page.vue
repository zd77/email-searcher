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

<script setup lang="ts">
    import { ref } from 'vue'
    import Pagination from './Pagination.vue'
    import Searchbar from './Searchbar.vue'
    import Table from './Table.vue'

    const VITE_BASE_URL = import.meta.env.VITE_BASE_URL
    const SEARCH_MAX_RESULT = 50

    type Items = {
        bodyMsg : string
        subject : string
        from    : string
        date    : string
        to      : string
    }

    type ApiResponse = {
        hits: {
            total: { value: number };
            hits: Array<{_source: Items}>
        }
    }

    const data = ref<Array<Items>>()
    const searchQuery = ref<string>('')
    const subheaderContent = ref<string>('')
    const bodyMessage = ref<string>('...')
    const currentPage = ref<number>(1)
    const totalEntries = ref<number>(0)
    const totalPages = ref<number>(1)


    const handleSearch = async ( queryString: string, page = 1 ): Promise<void> => {
        try {
            const credentials = btoa('admin:Complexpass#123')
            const response = await fetch(`${VITE_BASE_URL}/api/enron_emails/_search`,{
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
            const formattedResp: ApiResponse = await response.json()
            totalEntries.value = formattedResp.hits.total.value
            totalPages.value = Math.ceil(totalEntries.value/SEARCH_MAX_RESULT)
            const extractData = formattedResp.hits.hits
            data.value = extractData.map(({ _source }: {  _source: Items }) => ({
                to: _source.to,
                from: _source.from,
                subject: _source.subject,
                date: _source.date,
                bodyMsg: _source.bodyMsg,
            }))
        } catch( err ) {
            console.log("Error trying to fetch: ", err)
        }
    }

    const goToPage = ( page: number ) => {
        if( page >= 1 && page <= totalPages.value ) {
            currentPage.value = page
            handleSearch(searchQuery.value, page )
        }
    }

    const wordHighlighter = (query:string, text:string): string => {
        const regex = new RegExp(`\\b${query}\\b`, 'gi');
        const highlightedText = text.replace(regex, match => `<span class="font-bold">${match}</span>`)
        return highlightedText
    }

    const handleEnter = ( searchQ:string ): void => {
        searchQuery.value = searchQ
        bodyMessage.value = ' '
        subheaderContent.value = ' '
        currentPage.value = 1
        handleSearch(searchQuery.value)
    }

    const selectRow = ( selectedItem: Items ): void => {
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