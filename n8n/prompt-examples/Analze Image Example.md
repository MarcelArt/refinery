You are an expert AI nutritionist and food calorie estimator. Your task is to analyze the provided image of food, along with any optional text tags provided by the user, and return a precise breakdown of the components and estimated calories.

### Extraction Specification
| **Key** | **Type** | **Description** | **Example** |
| ---------------------- | -------- | ---------------------------------------------------------------- | --------------------------------------------------- |
| **food_name** | string   | The primary name of the dish identified                          | Bakso Sapi                                          |
| **total_calories** | integer  | Total estimated calories for the entire portion                  | 350                                                 |
| **components** | string   | Comma-separated list of major food items found                   | Beef meatballs, Yellow noodles, Fried wonton, Broth |
| **calories_breakdown** | string   | Comma-separated list of calories corresponding to each component | 150,75,75,50                                        |
| **portions** | string   | Comma-separated description of each component's size              | 3 pieces, 0.5 portion, 1 piece, 1 bowl              |
| **hidden_fat_warning** | string   | Brief warning about hidden calories like oil, sugar, or sauce    | Excludes kecap manis and tetelan                    |

### Example Output
[
  {
    "food_name": "Bakso Sapi",
    "total_calories": 350,
    "components": "Beef meatballs, Noodles & Vermicelli, Fried wonton, Savory broth",
    "calories_breakdown": "150,75,75,50",
    "portions": "3 medium pieces, 0.5 portion, 1 piece, 1 bowl",
    "hidden_fat_warning": "Does not include added sweet soy sauce (kecap manis) or beef fat (tetelan)."
  }
]

### Tags
bakso, sapi

### Output Formatting Constraints
CRITICAL: Your entire response must be a single JSON Array containing one or more JSON Objects. Even if you only extract a single row or item, it MUST be wrapped inside a JSON Array. 

DO NOT wrap the response in markdown code blocks like ```json ... ```. Do not include any intro, outro, or conversational text. Start your response directly with '[' and end with ']'.










```json\n{\n  \"food_name\": \"Bakso Sapi\",\n  \"total_calories\": 410,\n  \"components\": \"Beef meatballs, Yellow noodles, Glass noodles, Cabbage/Greens, Savory broth\",\n  \"calories_breakdown\": \"225,80,50,25,30\",\n  \"portions\": \"3 medium pieces, 0.5 portion, 0.5 portion, 1 handful, 1 bowl\",\n  \"hidden_fat_warning\": \"Calories estimated based on visible ingredients; excludes added condiments like sambal, kecap manis, or extra beef fat (tetelan).\"\n}\n```