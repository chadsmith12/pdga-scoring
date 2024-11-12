import pandas as pd
import matplotlib.pyplot as plt

# Load the CSV data into a DataFrame
data = pd.read_csv("1-team.csv")

# List of scoring types to consider
scoring_types = [
    "Event Winner", "Podiums", "Top 10s", "Hot Rounds", "Birdies",
    "Eagles Or Better", "Bogeys", "Double or Worse"
]

# Calculate the mean score for each scoring type
average_scores = data[scoring_types].mean()

# Plot the average contributions of each scoring type across all tournaments
plt.figure(figsize=(12, 6))
average_scores.plot(kind='bar', color='skyblue', edgecolor='black')
plt.title("Average Contribution of Each Scoring Type Across Tournaments")
plt.xlabel("Scoring Type")
plt.ylabel("Average Score")
plt.xticks(rotation=45)
plt.tight_layout()

# Save the plot as a single image file
plt.savefig("average_contributions-1.png")
plt.close()  # Close the figure to free memory
