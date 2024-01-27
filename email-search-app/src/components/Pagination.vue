<template>
    <div class="flex flex-row justify-between items-center mb-2">
        <button @click="goToPage(currentPage - 1)" class="flex items-center justify-center px-3 mr-2 h-8 text-sm font-medium text-white bg-gray-800 rounded hover:bg-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
            Prev
        </button>
        <span class="text-sm text-gray-700 dark:text-gray-400">
            Showing <span class="font-semibold text-gray-900 dark:text-white">{{ currentPage }}</span> to <span class="font-semibold text-gray-900 dark:text-white">{{ totalPages }} pages</span> of <span class="font-semibold text-gray-900 dark:text-white">{{ totalEntries }}</span> Entries
        </span>
        <button @click="goToPage(currentPage + 1)" class="flex items-center justify-center px-3 ml-2 h-8 text-sm font-medium text-white bg-gray-800 border-0 border-s border-gray-700 rounded hover:bg-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
            Next
        </button>
    </div>
</template>

<script setup lang="ts">
    const { currentPage, totalPages, goToPage, totalEntries } = defineProps(['currentPage', 'totalPages', 'goToPage', 'totalEntries'])
</script>