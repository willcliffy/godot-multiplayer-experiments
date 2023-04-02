from datetime import date, timedelta


if __name__ == "__main__":
    quarters = [1, 2, 3, 4]
    month_lengths = [4, 4, 4, 1]
    current_day = date(date.today().year, 1, 1)
    week = timedelta(days=7)
    current_month = 1
    for q in quarters:
        print(f"Quarter {q}")
        for month in month_lengths:
            current_week = 1
            print(f"\tMonth {current_month} has {month} weeks")
            current_month += 1
            for i in range(month):
                print(f"\t\tWeek {current_week} - {current_day}")
                current_week += 1
                current_day += week
