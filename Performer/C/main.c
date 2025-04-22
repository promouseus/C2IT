//
//  main.c
//  C2IT
//
//  Created by 00010110 B.V. on 15/03/2025.
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include <unistd.h>

#define MAX_LINE_LENGTH 1024  // Safe buffer size

// Function to read and print the file safely
void read_file_safely(const char *filename) {
    if (filename == NULL) {
        fprintf(stderr, "Error: NULL filename provided.\n");
        return;
    }

    FILE *file = fopen(filename, "r");
    if (file == NULL) {
        perror("Error opening file");
        return;
    }

    char buffer[MAX_LINE_LENGTH];

    printf("Reading file contents:\n");
    while (fgets(buffer, sizeof(buffer), file) != NULL) {
        size_t len = strnlen(buffer, sizeof(buffer));  // Safe length check

        // Check for truncation (line too long)
        if (len == MAX_LINE_LENGTH - 1 && buffer[len - 1] != '\n') {
            fprintf(stderr, "Warning: Line may be truncated.\n");
            int ch;
            while ((ch = fgetc(file)) != '\n' && ch != EOF) {
                // Consume remaining characters to avoid corruption
            }
        }

        printf("%s", buffer);  // Print each line
    }

    if (ferror(file)) {
        perror("Error reading file");
    }

    fclose(file);
}


// Function to update file safely
void update_file_safely(const char *filename, const char *search, const char *replace) {
    if (filename == NULL || search == NULL || replace == NULL) {
        fprintf(stderr, "Error: NULL parameter provided.\n");
        return;
    }

    FILE *file = fopen(filename, "r");
    if (file == NULL) {
        perror("Error opening file");
        return;
    }

    char temp_filename[FILENAME_MAX];
    snprintf(temp_filename, sizeof(temp_filename), "%s.tmp", filename);

    FILE *temp_file = fopen(temp_filename, "w");
    if (temp_file == NULL) {
        perror("Error opening temporary file");
        fclose(file);
        return;
    }

    char buffer[MAX_LINE_LENGTH];
    int modified = 0;  // Track if any modification is made

    while (fgets(buffer, sizeof(buffer), file) != NULL) {
        if (strstr(buffer, search) != NULL) {  // Found line to update
            snprintf(buffer, sizeof(buffer), "%s\n", replace);  // Replace safely
            modified = 1;  // Mark file as modified
        }

        fputs(buffer, temp_file);
    }

    if (ferror(file)) {
        perror("Error reading file");
    }

    fclose(file);
    fclose(temp_file);

    if (modified) {
        if (rename(temp_filename, filename) != 0) {  // Atomic replace
            perror("Error replacing file");
            remove(temp_filename);  // Clean up temp file if rename fails
        } else {
            printf("File updated successfully!\n");
        }
    } else {
        remove(temp_filename);  // No changes, delete temp file
        printf("No changes needed.\n");
    }
}

void print_current_directory(void) {
    char buffer[PATH_MAX];  // Use PATH_MAX for buffer size

    if (getcwd(buffer, sizeof(buffer)) == NULL) {
        perror("Error getting current directory");
        return;
    }

    printf("Current directory: %s\n", buffer);
}


int main(void) {
    const char *filename = "/Users/jeffrey/Documents/0 Signal/GitHub/C2IT/C2IT/C2IT/example.conf";
    
    print_current_directory();
    read_file_safely("/Users/jeffrey/Documents/0 Signal/GitHub/C2IT/C2IT/C2IT/example.conf");
    
    // First, read and analyze the file
    read_file_safely(filename);

    // Then, update the file if necessary
    // update_file_safely(filename, "old text", "new text");

    return 0;
}
