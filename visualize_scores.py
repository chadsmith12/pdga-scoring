import os
import sys
import pandas as pd
import matplotlib.pyplot as plt

def generate_visualizations(directory):
    # Check if the provided path is a directory
    if not os.path.isdir(directory):
        print(f"Error: {directory} is not a valid directory.")
        return

    # List of scoring types to consider (customize as needed)
    scoring_types = [
        "Event Winner", "Podiums", "Top 10s", "Hot Rounds", "Birdies",
        "Eagles Or Better", "Bogeys", "Double or Worse"
    ]

    # Process each CSV file in the provided directory
    for filename in os.listdir(directory):
        if filename.endswith(".csv"):
            csv_path = os.path.join(directory, filename)
            # Load the CSV data into a DataFrame
            data = pd.read_csv(csv_path)

            # Check if the required columns are present in the CSV
            columns_present = [col for col in scoring_types if col in data.columns]
            if not columns_present:
                print(f"No valid scoring columns found in {filename}. Skipping this file.")
                continue

            # Calculate the mean score for each scoring type
            average_scores = data[columns_present].mean()

            # Plot the average contributions of each scoring type
            plt.figure(figsize=(12, 6))
            average_scores.plot(kind='bar', color='skyblue', edgecolor='black')
            plt.title(f"Average Contribution of Each Scoring Type - {filename}")
            plt.xlabel("Scoring Type")
            plt.ylabel("Average Score")
            plt.xticks(rotation=45)
            plt.tight_layout()

            # Save the plot as an image with the same name as the CSV file
            image_filename = filename.replace(".csv", ".png")
            image_path = os.path.join(directory, image_filename)
            plt.savefig(image_path)
            plt.close()  # Close the figure to free memory

            print(f"Visualization saved for {filename} as {image_filename}")

# Run the function if this script is called directly
if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 generate_visualizations.py <directory_path>")
    else:
        directory_path = sys.argv[1]
        generate_visualizations(directory_path)
