export interface Category {
  id: number;
  name: string;
  description: string;
}

export const categories: Category[] = [
  { id: 1, name: 'Супы', description: 'Горячие первые блюда' },
  { id: 2, name: 'Салаты', description: 'Свежие и питательные салаты' },
  { id: 3, name: 'Горячие блюда', description: 'Основные блюда из мяса, рыбы и овощей' },
  { id: 4, name: 'Гарниры', description: 'Дополнения к основным блюдам' },
  { id: 5, name: 'Десерты', description: 'Сладкие блюда и выпечка' },
  { id: 6, name: 'Напитки', description: 'Безалкогольные напитки' },
  { id: 7, name: 'Завтраки', description: 'Блюда для утреннего меню' },
  { id: 8, name: 'Пасты', description: 'Итальянская паста и соусы' },
  { id: 9, name: 'Пицца', description: 'Разнообразные пиццы' },
  { id: 10, name: 'Соусы', description: 'Дополнительные соусы к блюдам' },
]; 