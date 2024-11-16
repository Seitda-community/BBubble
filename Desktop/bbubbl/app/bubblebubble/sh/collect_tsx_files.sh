find . \( -name "*.ts" -o -name "*.tsx" \) -type f -not -path "./node_modules/*" -not -path "./dist/*" -not -path "./coverage/*" -not -path "./pnpm-lock.yaml" -print0 | while IFS= read -r -d '' file; 
do
  echo "File: $file"
  echo "Content:"

  cat "$file"
  echo
done > ts_files_contents.txt

code ts_files_contents.txt

# copy to clipboard
pbcopy < ts_files_contents.txt