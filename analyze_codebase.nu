#!/opt/homebrew/bin/nu

# Analyze codebase by file type
def main [directory: string = "."] {
    print $"📊 Analyzing: ($directory)\n"
    
    cd $directory
    let files = (glob **/* | where ($it | path type) == file | wrap name)
    
    if ($files | is-empty) {
        print "No files found"
        return
    }
    
    let by_extension = ($files 
        | insert ext { |row| 
            let parsed = ($row.name | path parse)
            if ($parsed.extension | is-empty) {
                "no-ext"
            } else {
                $parsed.extension
            }
        }
        | group-by ext
        | transpose extension data
        | insert file_count { |row| $row.data | length }
        | insert total_lines { |row| 
            $row.data 
            | each { |file| 
                try { 
                    open $file.name | lines | length 
                } catch { 
                    0 
                }
            } 
            | math sum
        }
        | select extension file_count total_lines
        | sort-by total_lines -r
    )
    
    print "======================================================================"
    print ($by_extension)
    print ""
    print $"Total Files: ($by_extension | get file_count | math sum)"
    print $"Total Lines: ($by_extension | get total_lines | math sum)"
}
