#!/usr/bin/env python3
"""
Скрипт для импорта тестовых мест интереса в базу данных
"""

import psycopg2
from psycopg2.extras import Json
import os

# Тестовые данные POI для Москвы (советская эпоха)
SAMPLE_POIS = [
    {
        "name": "ВДНХ",
        "description": "Выставка достижений народного хозяйства СССР - грандиозный выставочный комплекс",
        "latitude": 55.8304,
        "longitude": 37.6325,
        "epoch": "soviet",
        "category": "architecture",
        "importance": 10,
        "year_built": 1939,
        "architect": "Вячеслав Олтаржевский",
        "style": "Сталинский ампир",
        "photos": [
            "https://example.com/vdnh1.jpg",
            "https://example.com/vdnh2.jpg"
        ]
    },
    {
        "name": "Павильон 'Космос' на ВДНХ",
        "description": "Легендарный павильон, посвященный космическим достижениям СССР",
        "latitude": 55.8283,
        "longitude": 37.6308,
        "epoch": "soviet",
        "category": "history",
        "importance": 9,
        "year_built": 1939,
        "style": "Сталинский ампир",
        "photos": []
    },
    {
        "name": "Фонтан 'Дружба народов'",
        "description": "Символ единства советских республик с позолоченными скульптурами",
        "latitude": 55.8277,
        "longitude": 37.6319,
        "epoch": "soviet",
        "category": "art",
        "importance": 9,
        "year_built": 1954,
        "architect": "Константин Топуридзе",
        "photos": []
    },
    {
        "name": "Рабочий и колхозница",
        "description": "Монументальная скульптура работы Веры Мухиной - символ советской эпохи",
        "latitude": 55.8311,
        "longitude": 37.6278,
        "epoch": "soviet",
        "category": "art",
        "importance": 10,
        "year_built": 1937,
        "architect": "Вера Мухина",
        "style": "Социалистический реализм",
        "photos": []
    },
    {
        "name": "МГУ им. М.В. Ломоносова",
        "description": "Главное здание МГУ - одна из знаменитых сталинских высоток",
        "latitude": 55.7033,
        "longitude": 37.5297,
        "epoch": "soviet",
        "category": "architecture",
        "importance": 10,
        "year_built": 1953,
        "architect": "Лев Руднев",
        "style": "Сталинский ампир",
        "photos": []
    },
    {
        "name": "Останкинская телебашня",
        "description": "Символ советской телевизионной эпохи, одно из высочайших зданий мира",
        "latitude": 55.8194,
        "longitude": 37.6119,
        "epoch": "soviet",
        "category": "architecture",
        "importance": 9,
        "year_built": 1967,
        "architect": "Николай Никитин",
        "photos": []
    },
    {
        "name": "Парк Горького",
        "description": "Центральный парк культуры и отдыха имени Горького - место отдыха советских граждан",
        "latitude": 55.7304,
        "longitude": 37.6012,
        "epoch": "soviet",
        "category": "culture",
        "importance": 8,
        "year_built": 1928,
        "photos": []
    },
    {
        "name": "Гостиница 'Украина'",
        "description": "Одна из семи сталинских высоток, образец советской архитектуры",
        "latitude": 55.7526,
        "longitude": 37.5676,
        "epoch": "soviet",
        "category": "architecture",
        "importance": 9,
        "year_built": 1957,
        "architect": "Аркадий Мордвинов",
        "style": "Сталинский ампир",
        "photos": []
    },
    {
        "name": "Третьяковская галерея",
        "description": "Крупнейший музей русского искусства в мире",
        "latitude": 55.7415,
        "longitude": 37.6206,
        "epoch": "imperial",
        "category": "art",
        "importance": 10,
        "year_built": 1856,
        "photos": []
    },
    {
        "name": "Красная площадь",
        "description": "Главная площадь Москвы, сердце столицы",
        "latitude": 55.7539,
        "longitude": 37.6208,
        "epoch": "medieval",
        "category": "history",
        "importance": 10,
        "year_built": 1493,
        "photos": []
    },
    {
        "name": "Храм Василия Блаженного",
        "description": "Собор Покрова Пресвятой Богородицы на Рву - шедевр русской архитектуры",
        "latitude": 55.7525,
        "longitude": 37.6231,
        "epoch": "medieval",
        "category": "religion",
        "importance": 10,
        "year_built": 1561,
        "architect": "Постник Яковлев",
        "photos": []
    },
    {
        "name": "Кремль",
        "description": "Московский Кремль - древняя крепость в центре Москвы",
        "latitude": 55.7520,
        "longitude": 37.6175,
        "epoch": "medieval",
        "category": "history",
        "importance": 10,
        "year_built": 1482,
        "photos": []
    }
]

def connect_db():
    """Подключение к PostgreSQL"""
    db_url = os.getenv(
        "DATABASE_URL",
        "postgresql://postgres:postgres@localhost:30101/audioguid"
    )
    return psycopg2.connect(db_url)

def import_pois(conn):
    """Импорт POI в базу данных"""
    cursor = conn.cursor()
    
    inserted = 0
    skipped = 0
    
    for poi in SAMPLE_POIS:
        try:
            # Проверяем, существует ли уже такой POI
            cursor.execute(
                "SELECT id FROM pois WHERE name = %s",
                (poi["name"],)
            )
            if cursor.fetchone():
                print(f"⏭️  Пропущено (уже существует): {poi['name']}")
                skipped += 1
                continue
            
            # Вставляем новый POI
            cursor.execute("""
                INSERT INTO pois (
                    name, description, latitude, longitude, epoch, category,
                    importance, year_built, architect, style, photos, created_at, updated_at
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, NOW(), NOW())
            """, (
                poi["name"],
                poi["description"],
                poi["latitude"],
                poi["longitude"],
                poi["epoch"],
                poi["category"],
                poi["importance"],
                poi.get("year_built"),
                poi.get("architect"),
                poi.get("style"),
                poi.get("photos", [])
            ))
            
            print(f"✅ Добавлено: {poi['name']}")
            inserted += 1
            
        except Exception as e:
            print(f"❌ Ошибка при добавлении '{poi['name']}': {e}")
    
    conn.commit()
    cursor.close()
    
    return inserted, skipped

def main():
    print("🚀 Импорт тестовых мест интереса...")
    print("=" * 60)
    
    try:
        conn = connect_db()
        print("✅ Подключено к базе данных")
        
        inserted, skipped = import_pois(conn)
        
        print("=" * 60)
        print(f"✅ Импорт завершен!")
        print(f"   Добавлено: {inserted}")
        print(f"   Пропущено: {skipped}")
        print(f"   Всего POI: {inserted + skipped}")
        
        conn.close()
        
    except Exception as e:
        print(f"❌ Ошибка: {e}")
        return 1
    
    return 0

if __name__ == "__main__":
    exit(main())

