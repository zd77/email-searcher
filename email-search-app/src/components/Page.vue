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
    import { handleSearch, Items } from '../handlers/handleSearch.js'

    const data = ref<Array<Items>>()
    const searchQuery = ref<string>('')
    const subheaderContent = ref<string>('')
    const bodyMessage = ref<string>('...')
    const currentPage = ref<number>(1)
    const totalEntries = ref<number>(0)
    const totalPages = ref<number>(1)

    const goToPage = ( page: number ) => {
        if( page >= 1 && page <= totalPages.value ) {
            currentPage.value = page
            handleSearch(searchQuery.value, totalEntries, totalPages, page )
                .then(( result ) => {
                    data.value = result
                })
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
        handleSearch(searchQuery.value, totalEntries, totalPages)
            .then(( result ) => {
                data.value = result
            })
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