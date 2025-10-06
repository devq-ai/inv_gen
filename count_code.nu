#!/opt/homebrew/bin/nu

cd .
let files = (glob **/* | where ($it | path type) == file | wrap name)

let by_extension = ($files 
    | insert ext { |row| 
        let parsed = ($row.name | path parse)
        if ($parsed.extension | is-empty) {
            "no-ext"
        } else {
            $parsed.extension
        }
    }
    | where ext != "no-ext"
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

print ($by_extension)
print ""
print $"Total Files (excluding no-ext): ($by_extension | get file_count | math sum)"
print $"Total Lines (excluding no-ext): ($by_extension | get total_lines | math sum)"
