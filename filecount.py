import os

def get_folder_size(path):
    total_size = 0
    count =0
    for dirpath, dirnames, filenames in os.walk(path):
        for filename in filenames:
            file_path = os.path.join(dirpath, filename)
            count+=1
            total_size += os.path.getsize(file_path)
    return total_size, count

def main():
    folder_path = "iheartwatson.net"  # Replace with the actual path

    total_size_bytes , count= get_folder_size(folder_path)
    total_size_gb = total_size_bytes / (1024 * 1024 * 1024)  # Convert bytes to GB
    total_size_mb = total_size_bytes / (1024 * 1024)  # Convert bytes to MB

    print(f"Total number of files: {count}")
    print(f"Total size: {total_size_gb:.2f} GB")
    print(f"Total size: {total_size_mb:.2f} MB")

if __name__ == "__main__":
    main()
